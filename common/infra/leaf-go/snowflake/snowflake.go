package snowflake

import (
	"context"
	"strconv"
	"time"
)

func (c *Creator) GetId() (int64, bool) {
	if c.working.Load() {
		return int64(c.snowNode.Generate()), true
	} else {
		return 0, false
	}
}

func (c *Creator) GetIdWithContext(ctx context.Context) (int64, error) {
	var id int64
	var ok bool

	for {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			id, ok = c.GetId()
			if !ok {
				time.Sleep(time.Millisecond * 50)
				continue
			}
		}
		break
	}
	return id, ctx.Err()
}

func (c *Creator) GetIdWithTimeout(timeout time.Duration) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return c.GetIdWithContext(ctx)
}

// heartCheck 心跳，定时上报时钟到本地和etcd
func (c *Creator) heartCheck() {
	ch := time.NewTicker(time.Millisecond * 200)
	key := "IdCreatorForever/" + c.name + "/" + c.addr

	for !c.local.Load() {
		select {
		case <-ch.C:
			timeout, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
			resp, err := c.client.Get(timeout, key)
			//当etcd请求失效时将本地存储的时钟作为依据
			var t int64

			if err != nil {
				t = c.lastTime
			} else {
				t, err = strconv.ParseInt(string(resp.Kvs[0].Value), 10, 64)
				if err != nil {
					t = c.lastTime
				}

				if time.Now().UnixMilli()-t <= 0 {
					//小步长
					if t-time.Now().UnixMilli() <= 500 {
						c.working.Store(false)
						// 等待双倍时间
						time.Sleep(time.Second * (time.Duration(t-time.Now().UnixMilli()) / time.Millisecond) * 2)
						c.working.Store(true)
					} else {
						panic("clock failed")
					}
				}
			}
			_, _ = c.client.Put(timeout, key, strconv.FormatInt(time.Now().UnixMilli(), 10))
			c.lastTime = time.Now().UnixMilli()
			cancel()
		}
	}

	for range ch.C {
		t := c.lastTime

		if time.Now().UnixMilli()-t <= 0 {
			if t-time.Now().UnixMilli() <= 500 {
				c.working.Store(false)
				time.Sleep(time.Second * (time.Duration(t-time.Now().UnixMilli()) / time.Millisecond) * 2)
				c.working.Store(true)
			} else {
				panic("clock failed")
			}
		}
		c.lastTime = time.Now().UnixMilli()
	}
}
