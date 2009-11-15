package client

import (
	"os";
	"log";
	"strings";
	"connection";
	. "protocol";
	. "redis";
)

type Client struct {
	conn  connection.Endpoint;
}

func Connect () (c *Client, err os.Error){
	spec := connection.DefaultSpec();
	c, err = ConnectWithSpec(spec);
	return;
}

func ConnectWithSpec (spec *connection.Spec) (c *Client, err os.Error) {
	c = new(Client);
	c.conn, err = connection.OpenNew (spec);
	return;
}

/** ----------------- REDIS INTERFACE ------------- **/
/** ----------------- G E N E R A T E ------------- **/

func (c *Client) Set (key string, arg1 []byte) (Error) {
	keybytes := strings.Bytes(key);

	_, err := c.conn.ServiceRequest (&SET, keybytes, arg1);
	return err;
}
// Redis EXISTS command.
func (c *Client) Exists (key string) (result bool, err Error){
	log.Stdout("=> Type " + key);
	keybytes := strings.Bytes(key);

	var resp Response;
	resp, err = c.conn.ServiceRequest (&EXISTS, keybytes);
	if err == nil { result = resp.GetBooleanValue(); }
	return result, err;
}

func (c *Client) Type (key string) (result string, err Error){
	log.Stdout("=> Type " + key);
	keybytes := strings.Bytes(key);

	var resp Response;
	resp, err = c.conn.ServiceRequest (&TYPE, keybytes);
	if err == nil { result = resp.GetStringValue(); }
	return result, err;
}

// Redis INCR command.
func (c *Client) Incr (key string) (num int64, err Error){
	log.Stdout("=> Incr " + key);
	keybytes := strings.Bytes(key);

	var resp Response;
	resp, err = c.conn.ServiceRequest (&INCR, keybytes);
	if err == nil { num = resp.GetNumberValue(); }
	return num, err;
}

func (c *Client) Ping() (Error){

	_, err := c.conn.ServiceRequest (&PING);
	return err;
}

func (c *Client) Get(key string) (data []byte, e Error) {
	log.Stdout("=> Get " + key);
	keybytes := strings.Bytes(key);
	
	var resp Response;
	resp, e = c.conn.ServiceRequest (&GET, keybytes);
	if e == nil { data = resp.GetBulkData(); }
	
	return data, e;
}
func (c *Client) Smembers (key string) (data [][]byte, e Error) {
	log.Stdout("=> SMEMBERS " + key);
	keybytes := strings.Bytes(key);
	var resp Response;
	resp, e = c.conn.ServiceRequest (&SMEMBERS, keybytes);
	if e == nil {
		data = resp.GetMultiBulkData();
	}
	return data, e;
}
