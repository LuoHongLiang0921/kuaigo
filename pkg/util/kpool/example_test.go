// @Description

package kpool

var (
	pool *Pool
)

//
//func TestPool(t *testing.T) {
//	pool := &pool{
//		Name: "pooling_test",
//		Dial: func() (Conn, error) {
//			return redis.BoolCmd{}
//		},
//		MaxIdle:         100,
//		MaxActive:       50,
//		IdleTimeout:     time.Duration(60) * time.Second,
//		MaxConnLifetime: time.Duration(120) * time.Second,
//	}
//
//	for index := 0; index < 10; index++ {
//		go func() {
//			conn, err := pool.Get(context.TODO())
//			if err != nil {
//				pool.Put(conn, true)
//				return
//			}
//			defer pool.Put(conn, false)
//			fmt.Println("IdleCount,", pool.IdleCount(), "ActiveCount,", pool.ActiveCount())
//			time.Sleep(10 * time.Second)
//		}()
//	}
//	time.Sleep(150 * time.Second)
//}
