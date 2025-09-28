package nundb

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/kqnd/nun-db-go/response"
)

type Client struct {
	Username string
	Database string

	// private
	responseHandler response.Handler
	conn            *websocket.Conn
	watchers        map[string][]func(interface{})
	pendings        []chan interface{}
	queue           sync.Mutex

	strictQueueWatchers bool
}

func NewClient(serverUrl, username, password string) (*Client, error) {
	watchers := make(map[string][]func(interface{}))
	pendings := make([]chan interface{}, 0)

	client := &Client{
		Username:            username,
		watchers:            watchers,
		pendings:            pendings,
		strictQueueWatchers: false,
	}

	u, err := url.Parse(serverUrl)
	if err != nil {
		return nil, err
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	handler := response.CreateHandler(&client.watchers, &client.pendings)
	client.responseHandler = handler
	client.conn = conn

	authMsg := fmt.Sprintf("auth %s %s", username, password)
	client.SendCommand(authMsg)

	go client.listen()

	return client, nil
}

func (c *Client) SetWatchersQueueMode(strict bool) {
	c.strictQueueWatchers = strict
}

func (c *Client) listen() {
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			return
		}

		c.responseHandler.SetEntireMessage(string(msg))
		c.responseHandler.SetPayload(string(msg))
		c.responseHandler.GettingValues()
		c.responseHandler.WatchingValues(c.strictQueueWatchers)
		c.responseHandler.AllDatabases()
		c.responseHandler.NoDBSelected()
		c.responseHandler.InvalidAuth()
		c.responseHandler.Keys()

	}
}

func (c *Client) verifyConnection() {
	if c.conn == nil {
		fmt.Println("no connection established with nundb server")
		os.Exit(1)
	}
}

func (c *Client) SendCommand(command string) error {
	c.verifyConnection()
	err := c.conn.WriteMessage(websocket.TextMessage, []byte(command))
	if err != nil {
		log.Fatalf("error during command: %s - error: %s", command, err)
		return err
	}
	return nil
}

func (c *Client) CreateDatabase(name, pwd string) error {
	if c.Database != "" {
		dbs, err := c.GetAllDatabases()
		if err != nil {
			fmt.Println("error fetching dbs:", err)
			return err
		}
		for _, db := range dbs {
			if db == name {
				return nil
			}
		}
	}
	c.SendCommand(fmt.Sprintf("create-db %s %s", name, pwd))
	return nil
}

func (c *Client) UseDatabase(name, pwd string) error {
	c.SendCommand(fmt.Sprintf("use-db %s %s", name, pwd))
	return nil
}

// a asd
func (c *Client) Set(key string, value interface{}) error {
	var strValue string
	switch v := value.(type) {
	case string:
		strValue = v
	default:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("error on parsing value to json: %w", err)
		}
		strValue = string(jsonBytes)
	}
	c.SendCommand(fmt.Sprintf("set %s %s", key, strValue))
	return nil
}

func (c *Client) Increment(key string, value int) {
	c.SendCommand(fmt.Sprintf("increment %s %d", key, value))
}

func (c *Client) Remove(key string) error {
	c.SendCommand("remove " + key)
	return nil
}

func (c *Client) RemoveAllWatchers() error {
	c.SendCommand("unwatch-all")
	c.watchers = make(map[string][]func(interface{}))
	return nil
}

func (c *Client) RemoveWatcher(key string) error {
	c.SendCommand("unwatch " + key)
	delete(c.watchers, key)
	return nil
}

func (c *Client) Watch(key string, cb func(interface{})) error {
	c.queue.Lock()
	c.watchers[key] = append(c.watchers[key], cb)
	c.queue.Unlock()
	c.SendCommand("watch " + key)
	return nil
}

func (c *Client) Get(key string) (interface{}, error) {
	ch := make(chan interface{})

	c.queue.Lock()
	c.pendings = append(c.pendings, ch)
	c.queue.Unlock()

	c.SendCommand("get " + key)
	return <-ch, nil
}

func (c *Client) GetAllDatabases() ([]string, error) {
	ch := make(chan interface{})

	c.queue.Lock()
	c.pendings = append(c.pendings, ch)
	c.queue.Unlock()

	c.SendCommand("debug list-dbs")

	res := <-ch
	dbs, ok := res.([]string)
	if !ok {
		return nil, fmt.Errorf("dbs is not a array")
	}
	return dbs, nil
}

func (c *Client) GetAllKeys() ([]string, error) {
	ch := make(chan interface{})

	c.queue.Lock()
	c.pendings = append(c.pendings, ch)
	c.queue.Unlock()

	c.SendCommand("keys")

	res := <-ch
	keys, ok := res.([]string)
	if !ok {
		return nil, fmt.Errorf("keys is not a array")
	}

	return keys, nil
}

func (c *Client) GetKeysStartingWith(prefix string) (interface{}, error) {
	ch := make(chan interface{})

	c.queue.Lock()
	c.pendings = append(c.pendings, ch)
	c.queue.Unlock()

	c.SendCommand("keys " + prefix)

	res := <-ch
	keys, ok := res.([]string)
	if !ok {
		return nil, fmt.Errorf("keys is not a array")
	}

	return keys, nil
}

func (c *Client) GetKeysEndingWith(suffix string) (interface{}, error) {
	ch := make(chan interface{})

	c.queue.Lock()
	c.pendings = append(c.pendings, ch)
	c.queue.Unlock()

	c.SendCommand("keys *" + suffix)

	res := <-ch
	keys, ok := res.([]string)
	if !ok {
		return nil, fmt.Errorf("keys is not a array")
	}

	return keys, nil
}
