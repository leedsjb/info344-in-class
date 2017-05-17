"use strict";

const express = require('express');
const Task = require('../models/tasks/task');

//export a function from this module 
//that accepts a tasks store implementation
//function supplies dependency injection (store)
module.exports = function(store) {
  
    // create a new mux, similar to GoLang
    let router = express.Router();

    // must have async before function to use await
    router.get('/v1/tasks', async (req, res, next) => { // next provided centralized error handling

        try{
            // throw new Error('testing');

            let tasks = await store.getAll(); // await makes this call synchronous
            res.json(tasks);
        } catch(err){
            next(err); // centralized error handling
        }

        // "old style (except not)" of using a promise
        store.getAll()
            .then(tasks => {
                res.json(tasks);
            })
            .catch(next);

    });

    router.post('/v1/tasks', async (req, res, next) => {

        try{
            let task = new Task(req.body) // uses body parser to define request body

            let err = task.validate(); // validate newly created task

            if (err){
                res.status(400).send(err.message);
            } else {
                let result = await store.insert(task);
                res.json(task);
            }
        }catch(err){
            next(err);
        }
    })

    router.patch('/v1/tasks/:taskID', async (req, res, next)=> {
        let taskID = req.params.taskID;

        try {
            let updatedTask = store.setComplete(taskID, req.body.complete);
            res.json(updatedTask);
        } catch(err) {
            next(err);
        }
    })

    router.delete('/v1/tasks/:taskID', async (req, res, next) => {
        let taskID = req.params.taskID;

        try{
            await store.delete(taskID);
            res.send(`deleted task ${taskID}`);
        } catch (err){
            next(err);
        }
    })

    return router;
};
