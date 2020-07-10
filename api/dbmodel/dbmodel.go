package dbmodel

import (
	"errors"
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const MONGO_CONNSTRING = "mongo:27017"

// Data model for our document collection
type Document struct {
	Id       bson.ObjectId `bson:"_id"`
	ShortUrl string        `bson:"shorturl"`
	LongUrl  string        `bson:"longurl"`
}

type Connection struct {
	originalSession *mgo.Session
}

func NewConnection() (conn *Connection) {
	conn = new(Connection)
	conn.createConnection()
	return
}

func (c *Connection) createConnection() (err error) {
	fmt.Println("Connecting to mongo....")
	c.originalSession, err = mgo.Dial(MONGO_CONNSTRING)
	if err == nil {
		fmt.Println("Connection established to mongo")
		collection := c.originalSession.DB("ShortifyDB").C("Urls")
		if collection == nil {
			err = errors.New("Collection could not be created, maybe need to create it manually")
		}
		// Create an unique index so shorturls do not repeat
		// Key is shorturl
		// Unique is true
		index := mgo.Index{
			Key:      []string{"$text:shorturl"},
			Unique:   true,
			DropDups: true,
		}
		collection.EnsureIndex(index)
	} else {
		fmt.Printf("Error occured while creating mongodb connection: %s", err.Error())
	}
	return
}

func (c *Connection) CloseConnection() {
	if c.originalSession != nil {
		c.originalSession.Close()
	}
}

func (c *Connection) getSessionAndCollection() (session *mgo.Session, collection *mgo.Collection, err error) {
	if c.originalSession != nil {
		session = c.originalSession.Copy()
		collection = session.DB("ShortifyDB").C("Urls")
	} else {
		err = errors.New("No original session found")
	}
	return
}

func (c *Connection) FindShortURL(long string) (short string, err error) {
	result := Document{}
	session, collection, err := c.getSessionAndCollection()
	if err != nil {
		return
	}
	defer session.Close()
	err = collection.Find(bson.M{"longurl": long}).One(&result)
	if err != nil {
		return
	}
	return result.ShortUrl, nil
}

func (c *Connection) FindLongURL(short string) (long string, err error) {
	result := Document{}
	session, collection, err := c.getSessionAndCollection()
	if err != nil {
		return
	}
	defer session.Close()
	err = collection.Find(bson.M{"shorturl": short}).One(&result)
	if err != nil {
		return
	}
	return result.LongUrl, nil
}

func (c *Connection) AddURL(longUrl string, shortUrl string) (err error) {
	session, collection, err := c.getSessionAndCollection()
	if err == nil {
		defer session.Close()
		err = collection.Insert(
			&Document{
				Id:       bson.NewObjectId(),
				ShortUrl: shortUrl,
				LongUrl:  longUrl,
			},
		)
		if err != nil {
			// Unlikely, but still want to check if too many urls created and one collides
			if mgo.IsDup(err) {
				err = errors.New("Duplicate name exists for the shorturl")
			}
		}
	}
	return
}
