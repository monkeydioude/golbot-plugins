package reddithot

import (
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
)

type subList struct {
	cachePath string
}

func newSubList(cacheDir string) *subList {
	return &subList{
		cachePath: fmt.Sprintf("%s%s", cacheDir, "sublistByDiscordID"),
	}
}

func (s *subList) addSavedSubFromCache(r *redditHot) error {
	db, err := bolt.Open(s.cachePath, 0666, nil)
	defer db.Close()
	if err != nil {
		log.Printf("[ERR ] %s\n", err.Error())
		return err
	}
	err = db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			channelID := string(name)
			return b.ForEach(func(key, _ []byte) error {
				sub := string(key)
				watchSub(channelID, string(sub), r)

				return nil
			})
		})
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *subList) addSubToSubList(channelID, sub string) error {
	db, err := bolt.Open(s.cachePath, 0666, nil)
	defer db.Close()
	if err != nil {
		log.Printf("[ERR ] %s\n", err)
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(channelID))

		if b == nil {
			b, err = tx.CreateBucket([]byte(channelID))

			if err != nil {
				return err
			}
		}
		return b.Put([]byte(sub), []byte("1"))
	})
	if err != nil {
		log.Printf("[ERR ] %s\n", err)
		return err
	}
	return nil
}

func (s *subList) removeSubFromSubList(channelID, sub string) error {
	db, err := bolt.Open(s.cachePath, 0666, nil)
	defer db.Close()
	if err != nil {
		log.Printf("[ERR ] %s\n", err)
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(channelID))

		if b == nil {
			return nil
		}
		return b.Delete([]byte(sub))
	})
	if err != nil {
		log.Printf("[ERR ] %s\n", err)
		return err
	}
	return nil
}
