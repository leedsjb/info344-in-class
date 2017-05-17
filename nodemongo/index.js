"use strict";

const express = require('express');
const morgan = require('morgan');
const cors = require('cors');
const bodyParser = require('body-parser');
const mongodb = require('mongodb');
const TaskStore = require('./models/tasks/mongostore');

const port = process.env.PORT || 80;
const host = process.env.HOST || '';
const mongoAddr = process.env.MONGOADDR || 'localhost:27017'; // changed from 27017 to 27018 due to other running container

//create an Experss application
const app = express();
//add request logging
app.use(morgan(process.env.LOGFORMAT || 'dev'));
//add CORS headers
app.use(cors());
//add middleware that parses
//any JSON posted to this app.
//the parsed data will be available
//on the req.body property
app.use(bodyParser.json());

//TODO: connect to the Mongo database
//add the tasks handlers
//and start listening for HTTP requests

mongodb.MongoClient.connect(`mongodb://${mongoAddr}/demo`) // demo is the database we want to write to
    .then(db => {

        let colTasks = db.collection('tasks'); // set collection name for demo database
        let store = new TaskStore(colTasks);
        let handlers = require('./handlers/tasks.js');
        app.use(handlers(store));

        // error handler
        app.use((err, req, res, next)=> {
            console.error(err);
            res.status(500).send(err.message);
        });

        app.listen(port, host, ()=> {
            console.log(`server is listening at http://${host}:${port}...`);
        } )

    })
    .catch(err => {
        console.error(err);
    })
