"use strict";

//TODO: implement a Task class
//for the task model, and export
//it from this module

// JavaScript ES6 syntax
class Task {
    constructor(props){
        Object.assign(this, props); // copies fields from props to this
    }

    validate(){
        if(!this.title) { // falsey, does not exist or empty string
            return new Error('you must supply a tile');
        }
    }

}

module.exports = Task; // makes Task available in other js files, must use module.exports.foo to export multiple classes