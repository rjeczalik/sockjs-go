package sockjs

import (
	"crypto/rand"
	"encoding/hex"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var DefaultClient = &Client{}

type Client struct {
	Client    *http.Client
	ReadSize  int
	WriteSize int
	Timeout   time.Duration
}

func (c *Client) timeout() time.Duration {
	if c.Timeout != 0 {
		return c.Timeout
	}
	return 30 * time.Second
}

func (c *Client) client() *http.Client {
	if c.Client != nil {
		return c.Client
	}
	return http.DefaultClient
}

func (c *Client) dialer() *websocket.Dialer {
	dialer := &websocket.Dialer{
		ReadBufferSize:   c.ReadSize,
		WriteBufferSize:  c.WriteSize,
		HandshakeTimeout: c.Timeout,
	}

	if t, ok := c.client().Transport.(*http.Transport); ok {
		dialer.NetDial = t.Dial
	}

	if dialer.NetDial == nil {
		dialer.NetDial = (&net.Dialer{
			Timeout:   c.timeout(),
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).Dial
	}

	return dialer
}

func (c *Client) Dial(uri string) (Session, error) {

	return nil, nil
}

func (c *Client) DialWebsocket(uri string) (Session, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	h := http.Header{
		"Origin": {u.Scheme + "://" + u.Host},
	}

	if u.Scheme == "https" {
		u.Scheme = "wss"
	} else {
		u.Scheme = "ws"
	}

	if _, _, err := net.SplitHostPort(u.Host); err != nil {
		if u.Scheme == "wss" {
			u.Host = net.JoinHostPort(u.Host, "443")
		} else {
			u.Host = net.JoinHostPort(u.Host, "80")
		}
	}

	if strings.HasSuffix(u.Path, "/") {
		u.Path = u.Path + "/"
	}

	u.Path = u.Path + serverID + "/" + sessionID + "/websocket"

	serverID := randShortID()
	sessionID := randString(20)

	conn, _, err := c.dialer().Dial(u.String(), h)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		URL: u,
	}

	sess := newSession(req, sessionID, c.disconnectDelay(), c.heartbeatDelay())
}

func Dial(url string) (Session, error) { return DefaultClient.Dial(url) }

func randShortID() string {
	p := make([]byte, 3)
	rand.Read(p)
	p0, p1, p2 := int64(p[0]), int64(p[1]), int64(p[2])
	return strconv.FormatInt(100+100*(p0%10)+10*(p1%10)+(p2%10), 10)
}

func randString(n int) string {
	p := make([]byte, n/2+1)
	rand.Read(p)
	return hex.EncodeToString(p)[:n]
}
