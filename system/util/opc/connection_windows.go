// +build windows

package opc

import (
	"errors"
	"fmt"
	"sync"
	"time"

	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func init() {
	OleInit()
}

//OleInit initializes OLE.
func OleInit() {
	ole.CoInitializeEx(0, 0)
}

//OleRelease realeses OLE resources in opcAutomation.
func OleRelease() {
	ole.CoUninitialize()
}

//AutomationObject loads the OPC Automation Wrapper and handles to connection to the OPC Server.
type AutomationObject struct {
	unknown *ole.IUnknown
	object  *ole.IDispatch
}

//CreateBrowser returns the OPCBrowser object from the OPCServer.
//It only works if there is a successful connection.
func (ao *AutomationObject) CreateBrowser() (*Tree, error) {
	// check if server is running, if not return error
	if !ao.IsConnected() {
		return nil, errors.New("Cannot create browser because we are not connected.")
	}

	// create browser
	browser, err := oleutil.CallMethod(ao.object, "CreateBrowser")
	if err != nil {
		return nil, errors.New("Failed to create OPCBrowser")
	}

	// move to root
	oleutil.MustCallMethod(browser.ToIDispatch(), "MoveToRoot")

	// create tree
	root := Tree{"root", nil, []*Tree{}, []Leaf{}}
	buildTree(browser.ToIDispatch(), &root)

	return &root, nil
}

//获取子节点内容
func (ao *AutomationObject) GetChildNode(node []string) (result []Leaf, err error) {

	defer func() {
		if err := recover(); err != nil {
			err = fmt.Errorf("目录异常")
		}
		return
	}()

	// check if server is running, if not return error
	if !ao.IsConnected() {
		return nil, errors.New("Cannot create browser because we are not connected.")
	}

	// create browser
	browser, err := oleutil.CallMethod(ao.object, "CreateBrowser")
	if err != nil {
		return nil, errors.New("Failed to create OPCBrowser")
	}
	oleutil.MustCallMethod(browser.ToIDispatch(), "MoveToRoot")
	if len(node) > 0 {
		err = fmt.Errorf("目录异常")
		oleutil.MustCallMethod(browser.ToIDispatch(), "MoveTo", node)
		err = nil
	}
	//  oleutil.MustCallMethod(browser.ToIDispatch(), "MoveDown", "test")

	oleutil.MustCallMethod(browser.ToIDispatch(), "ShowBranches").ToIDispatch()
	count := oleutil.MustGetProperty(browser.ToIDispatch(), "Count").Value().(int32)

	for i := 1; i <= int(count); i++ {
		item := oleutil.MustCallMethod(browser.ToIDispatch(), "Item", i).Value()

		tag := oleutil.MustCallMethod(browser.ToIDispatch(), "GetItemID", item).Value()

		//datetype := oleutil.MustGetProperty(browser, "DataType").Value()
		//datetype := oleutil.MustCallMethod(browser, "GetItemProperties", item).Value()
		var datetype int16 = 1
		l := Leaf{Name: item.(string), Tag: tag.(string), Type: datetype}
		result = append(result, l)
	}

	oleutil.MustCallMethod(browser.ToIDispatch(), "ShowLeafs")
	count = oleutil.MustGetProperty(browser.ToIDispatch(), "Count").Value().(int32)
	for i := 1; i <= int(count); i++ {

		item := oleutil.MustCallMethod(browser.ToIDispatch(), "Item", i).Value()
		tag := oleutil.MustCallMethod(browser.ToIDispatch(), "GetItemID", item).Value()

		//datetype := oleutil.MustGetProperty(browser, "DataType").Value()
		//datetype := oleutil.MustCallMethod(browser, "GetItemProperties", item).Value()
		var datetype int16 = 2
		l := Leaf{Name: item.(string), Tag: tag.(string), Type: datetype}
		result = append(result, l)

	}

	return
}

//buildTree runs through the OPCBrowser and creates a tree with the OPC tags
func buildTree(browser *ole.IDispatch, branch *Tree) {
	var count int32

	logger.Println("Entering branch:", branch.Name)

	// loop through leafs
	oleutil.MustCallMethod(browser, "ShowLeafs").ToIDispatch()
	count = oleutil.MustGetProperty(browser, "Count").Value().(int32)

	logger.Println("\tLeafs count:", count)

	for i := 1; i <= int(count); i++ {

		item := oleutil.MustCallMethod(browser, "Item", i).Value()
		tag := oleutil.MustCallMethod(browser, "GetItemID", item).Value()

		//datetype := oleutil.MustGetProperty(browser, "DataType").Value()
		//datetype := oleutil.MustCallMethod(browser, "GetItemProperties", item).Value()
		var datetype int16 = 0
		l := Leaf{Name: item.(string), Tag: tag.(string), Type: datetype}

		logger.Println("\t", i, l)

		branch.Leaves = append(branch.Leaves, l)
	}

	// loop through branches
	oleutil.MustCallMethod(browser, "ShowBranches").ToIDispatch()
	count = oleutil.MustGetProperty(browser, "Count").Value().(int32)

	logger.Println("\tBranches count:", count)

	for i := 1; i <= int(count); i++ {

		nextName := oleutil.MustCallMethod(browser, "Item", i).Value()

		logger.Println("\t", i, "next branch:", nextName)

		// move down
		oleutil.MustCallMethod(browser, "MoveDown", nextName)

		// recursively populate tree
		nextBranch := Tree{nextName.(string), branch, []*Tree{}, []Leaf{}}
		branch.Branches = append(branch.Branches, &nextBranch)
		buildTree(browser, &nextBranch)

		// move up and set branches again
		oleutil.MustCallMethod(browser, "MoveUp")
		oleutil.MustCallMethod(browser, "ShowBranches").ToIDispatch()
	}

	logger.Println("Exiting branch:", branch.Name)

}

//断开连接AutomationItems
func (ao *AutomationObject) Disconnect() error {

	// check if server is running, if yes then disconnect
	if ao.IsConnected() {
		_, err := oleutil.CallMethod(ao.object, "Disconnect")
		if err != nil {
			logger.Println("Failed to disconnect. Trying to connect anyway..")
		}
		return err
	}

	return nil
}

//Connect establishes a connection to the OPC Server on node.
//It returns a reference to AutomationItems and error message.
func (ao *AutomationObject) Connect(server string, node string) (*AutomationItems, error) {

	// check if server is running, if yes then disconnect
	if ao.IsConnected() {
		_, err := oleutil.CallMethod(ao.object, "Disconnect")
		if err != nil {
			logger.Println("Failed to disconnect. Trying to connect anyway..")
		}
	}

	// try to connect to opc server and check for error
	logger.Printf("Connecting to %s on node %s\n", server, node)
	_, err := oleutil.CallMethod(ao.object, "Connect", server, node)
	if err != nil {
		logger.Println("Connection failed.")
		return nil, errors.New("Connection failed")
	}

	// set up opc groups and items
	opcGroups, err := oleutil.GetProperty(ao.object, "OPCGroups")
	if err != nil {
		//logger.Println(err)
		return nil, errors.New("cannot get OPCGroups property")
	}
	opcGrp, err := oleutil.CallMethod(opcGroups.ToIDispatch(), "Add")
	if err != nil {
		// logger.Println(err)
		return nil, errors.New("cannot add new OPC Group")
	}
	addItemObject, err := oleutil.GetProperty(opcGrp.ToIDispatch(), "OPCItems")
	if err != nil {
		// logger.Println(err)
		return nil, errors.New("cannot get OPC Items")
	}

	opcGroups.ToIDispatch().Release()
	opcGrp.ToIDispatch().Release()

	logger.Println("Connected.")

	return NewAutomationItems(addItemObject.ToIDispatch()), nil
}

//TryConnect loops over the nodes array and tries to connect to any of the servers.
func (ao *AutomationObject) TryConnect(server string, nodes []string) (*AutomationItems, error) {
	var errResult string
	for _, node := range nodes {
		items, err := ao.Connect(server, node)
		if err == nil {
			return items, err
		}
		errResult = errResult + err.Error() + "\n"
	}
	return nil, errors.New("TryConnect was not successful: " + errResult)
}

//IsConnected check if the server is properly connected and up and running.
func (ao *AutomationObject) IsConnected() bool {
	if ao.object == nil {
		return false
	}
	stateVt, err := oleutil.GetProperty(ao.object, "ServerState")
	if err != nil {
		logger.Println("GetProperty call for ServerState failed", err)
		return false
	}
	if stateVt.Value().(int32) != OPCRunning {
		return false
	}
	return true
}

//GetOPCServers returns a list of Prog ID on the specified node
func (ao *AutomationObject) GetOPCServers(node string) []string {
	progids, err := oleutil.CallMethod(ao.object, "GetOPCServers", node)
	if err != nil {
		logger.Println("GetOPCServers call failed.")
		return []string{}
	}

	var servers_found []string
	for _, v := range progids.ToArray().ToStringArray() {
		if v != "" {
			servers_found = append(servers_found, v)
		}
	}
	return servers_found
}

//Close releases the OLE objects in the AutomationObject.
func (ao *AutomationObject) Close() {
	if ao.object != nil {
		ao.object.Release()
	}
	if ao.unknown != nil {
		ao.unknown.Release()
	}
}

//NewAutomationObject connects to the COM object based on available wrappers.
func NewAutomationObject() *AutomationObject {
	wrappers := []string{"OPC.Automation.1", "Graybox.OPC.DAWrapper.1"}
	//wrappers := []string{"OPC.Automation.1"}
	var err error
	var unknown *ole.IUnknown
	for _, wrapper := range wrappers {
		unknown, err = oleutil.CreateObject(wrapper)
		if err == nil {
			logger.Println("Loaded OPC Automation object with wrapper", wrapper)
			break
		}
		logger.Println("Could not load OPC Automation object with wrapper", wrapper)
	}
	if err != nil {
		return &AutomationObject{}
	}

	opc, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		fmt.Println("Could not QueryInterface")
		return &AutomationObject{}
	}
	object := AutomationObject{
		unknown: unknown,
		object:  opc,
	}
	return &object
}

//AutomationItems store the OPCItems from OPCGroup and does the bookkeeping
//for the individual OPC items. Tags can added, removed, and read.
type AutomationItems struct {
	addItemObject *ole.IDispatch
	items         map[string]*ole.IDispatch
}

//addSingle adds the tag and returns an error. Client handles are not implemented yet.
func (ai *AutomationItems) addSingle(tag string) (result_err error) {
	defer func() {
		if err := recover(); err != nil {
			result_err = fmt.Errorf("%s点位不存在", tag)
		}
	}()
	clientHandle := int32(1)
	item, err := oleutil.CallMethod(ai.addItemObject, "AddItem", tag, clientHandle)
	if err != nil {
		return errors.New(tag + ":" + err.Error())
	}
	ai.items[tag] = item.ToIDispatch()
	return
}

//Add accepts a variadic parameters of tags.
func (ai *AutomationItems) Add(tags ...string) map[string]error {

	var errResult map[string]error = make(map[string]error, 0)
	for _, tag := range tags {
		err := ai.addSingle(tag)
		errResult[tag] = err
	}
	return errResult
}

//Remove removes the tag.
func (ai *AutomationItems) Remove(tag string) {
	item, ok := ai.items[tag]
	if ok {
		item.Release()
	}
	delete(ai.items, tag)
}

/*
 * FIX:
 * some opc servers sometimes returns an int32 Quality, that produces panic
 */
func ensureInt16(q interface{}) int16 {
	if v16, ok := q.(int16); ok {
		return v16
	}
	if v32, ok := q.(int32); ok && v32 >= -32768 && v32 < 32768 {
		return int16(v32)
	}
	return 0
}

//readFromOPC reads from the server and returns an Item and error.
func (ai *AutomationItems) readFromOpc(opcitem *ole.IDispatch) (reItem Item, result_err error) {

	defer func() {
		if err := recover(); err != nil {
			result_err = fmt.Errorf("tag异常或服务异常")
		}

	}()

	v := ole.NewVariant(ole.VT_R4, 0)
	q := ole.NewVariant(ole.VT_INT, 0)
	ts := ole.NewVariant(ole.VT_DATE, 0)

	//read tag from opc server and monitor duration in seconds
	t := time.Now()
	_, err := oleutil.CallMethod(opcitem, "Read", OPCDevice, &v, &q, &ts) //这里控制是否为缓存
	opcReadsDuration.Observe(time.Since(t).Seconds())

	if err != nil {
		opcReadsCounter.WithLabelValues("failed").Inc()
		return Item{}, err
	}
	opcReadsCounter.WithLabelValues("success").Inc()

	reItem = Item{
		Value:     v.Value(),
		Quality:   ensureInt16(q.Value()), // FIX: ensure the quality value is int16
		Timestamp: ts.Value().(time.Time),
	}
	return
}

//writeToOPC writes value to opc tag and return an error
func (ai *AutomationItems) writeToOpc(opcitem *ole.IDispatch, value interface{}) (result_err error) {
	defer func() {
		if err := recover(); err != nil {
			result_err = fmt.Errorf("写入时点位发生异常")
		}
	}()

	_, err := oleutil.CallMethod(opcitem, "Write", value)
	if err != nil {
		// TODO: Prometheus Monitoring
		//opcWritesCounter.WithLabelValues("failed").Inc()
		return err
	}
	//opcWritesCounter.WithLabelValues("failed").Inc()
	return
}

//Close closes the OLE objects in AutomationItems.
func (ai *AutomationItems) Close() {
	if ai != nil {
		for key, opcitem := range ai.items {
			opcitem.Release()
			delete(ai.items, key)
		}
		ai.addItemObject.Release()
	}
}

//NewAutomationItems returns a new AutomationItems instance.
func NewAutomationItems(opcitems *ole.IDispatch) *AutomationItems {
	ai := AutomationItems{addItemObject: opcitems, items: make(map[string]*ole.IDispatch)}
	return &ai
}

//opcRealServer implements the Connection interface.
//It has the AutomationObject embedded for connecting to the server
//and an AutomationItems to facilitate the OPC items bookkeeping.
type opcConnectionImpl struct {
	*AutomationObject
	*AutomationItems
	Server string
	Nodes  []string
	mu     sync.Mutex
}

//ReadItem returns an Item for a specific tag.
func (conn *opcConnectionImpl) ReadItem(tag string) Item {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	opcitem, ok := conn.AutomationItems.items[tag]
	if ok {
		item, err := conn.AutomationItems.readFromOpc(opcitem)
		if err == nil {
			return item
		}
		logger.Printf("Cannot read %s: %s. Trying to fix.", tag, err)
		conn.fix()
	} else {
		logger.Printf("Tag %s not found. Add it first before reading it.", tag)
	}
	return Item{}
}

//Write writes a value to the OPC Server.
func (conn *opcConnectionImpl) Write(tag string, value interface{}) error {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	opcitem, ok := conn.AutomationItems.items[tag]
	if ok {
		return conn.AutomationItems.writeToOpc(opcitem, value)
	}
	logger.Printf("Tag %s not found. Add it first before writing to it.", tag)
	return errors.New("No Write performed")
}

//Read returns a map of the values of all added tags.
func (conn *opcConnectionImpl) Read() map[string]Item {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	allTags := make(map[string]Item)
	for tag, opcitem := range conn.AutomationItems.items {
		item, err := conn.AutomationItems.readFromOpc(opcitem)
		if err != nil {
			logger.Printf("Cannot read %s: %s. Trying to fix.", tag, err)
			conn.fix()
			continue
		}
		allTags[tag] = item
	}
	return allTags
}

//Tags returns the currently active tags
func (conn *opcConnectionImpl) Tags() []string {
	var tags []string
	if conn.AutomationItems != nil {
		for tag, _ := range conn.AutomationItems.items {
			tags = append(tags, tag)
		}
	}
	return tags

}

//fix tries to reconnect if connection is lost by creating a new connection
//with AutomationObject and creating a new AutomationItems instance.
func (conn *opcConnectionImpl) fix() {
	var err error
	if !conn.IsConnected() {
		for {
			tags := conn.Tags()
			conn.AutomationItems.Close()
			conn.AutomationItems, err = conn.TryConnect(conn.Server, conn.Nodes)
			if err != nil {
				logger.Println(err)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if conn.Add(tags...) == nil {
				logger.Printf("Added %d tags", len(tags))
			}
			break
		}
	}
}

//Close closes the embedded types.
func (conn *opcConnectionImpl) Close() {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if conn.AutomationObject != nil {
		conn.AutomationObject.Close()
		//这里加断开连接
		conn.AutomationObject.Disconnect()
	}
	if conn.AutomationItems != nil {
		conn.AutomationItems.Close()
	}

}

//NewConnection establishes a connection to the OpcServer object.
func NewConnection(server string, nodes []string, tags []string) (Connection, error) {
	object := NewAutomationObject()
	items, err := object.TryConnect(server, nodes)
	if err != nil {
		return &opcConnectionImpl{}, err
	}
	items.Add(tags...)

	conn := opcConnectionImpl{
		AutomationObject: object,
		AutomationItems:  items,
		Server:           server,
		Nodes:            nodes,
	}

	return &conn, nil
}

//CreateBrowser creates an opc browser representation
func CreateBrowser(server string, nodes []string) (*Tree, error) {
	object := NewAutomationObject()
	defer object.Close()
	_, err := object.TryConnect(server, nodes)
	if err != nil {
		return nil, err
	}
	return object.CreateBrowser()
}

func GetChildNode(server string, nodes []string, point []string) ([]Leaf, error) {
	object := NewAutomationObject()
	defer object.Close()
	_, err := object.TryConnect(server, nodes)
	if err != nil {
		return nil, err
	}
	return object.GetChildNode(point)
}
