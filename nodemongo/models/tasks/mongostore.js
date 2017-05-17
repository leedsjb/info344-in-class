/** 
 * mongostore.js
 * Created May 17, 2017 from Dr. Dave Stearns
 * Modified: 
*/

"use strict";

const mongodb = require('mongodb'); //for mongodb.ObjectID()

/**
 * MongoStore is a concrete store for Task models
 */
class MongoStore {
    /**
     * Constructs a new MongoStore
     * @param {mongodb.Collection} collection 
     */
    constructor(collection) {
        this.collection = collection;
    }

    /**
     * getAll returns all tasks in the store
     */
    getAll() {
        return this.collection.find().toArray() // creates and populates array simultaneously
    }

    /**
     * insert inserts a new Task into the store
     * @param {Task} task 
     */
    insert(task) {
        return this.collection.insert(task) // automatically checks and corrects for _id field mongo needs
    }

    /**
     * setComplete sets the complete status of the task
     * @param {string} id 
     * @param {bool} complete 
     */
    async setComplete(id, complete) { // ** note async, v. 7.10.0 and up only

        let options = {returnOriginal: false};
        let updates = {$set: {complete: complete}};

        let oid = new mongodb.ObjectID(id);

        let result = await this.collection.findOneAndUpdate({_id: oid},updates, options);

        return result.value;

    }

    /**
     * delete deletes the task with the given object ID
     * @param {string} id 
     */
    delete(id) {
        return this.collection.deleteOne({_id: new mongodb.ObjectID(id)});
    }
}

//export the class
module.exports = MongoStore;
