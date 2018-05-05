package main


func main() {
	cli := CLI{}
	cli.Run()
	//db, err := bolt.Open(dbFile,0600,nil)
	//if err != nil {
	//	log.Panic(err)
	//}
	//err = db.View(func(tx *bolt.Tx) error {
	//	b := tx.Bucket([]byte(blocksBucket))
	//	tip := b.Get([]byte("1"))
	//	fmt.Println(tip)
	//	return nil
	//})
}
