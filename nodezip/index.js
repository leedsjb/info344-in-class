"use strict"; // prevents creation of new var by type, among other compiler requirements

const express = require('express');
const cors = require('cors');
const morgan = require('morgan');
// require function can also import json files
const zips = require('./zips.json') // starting param w/ ./ tells node this is a relative filepath
console.log('loaded %d zips', zips.length);

// my try
// let zipCityIndex; // create javascript object
// create a js object w/ a key for each city where value is array of zipcodes
// iterate over json
// zips.array.forEach(function(element) {
//     element.toLowerCase()// lowercase element
//     zipCityIndex[element] = element.**zips**;
// }, this);

// stearn's code:
const zipCityIndex = zips.reduce((index,record) => {
    let cityLower = record.city.toLowerCase(); // convert city name to lowercase
    let zipsForCity = index[cityLower];
    if (!zipsForCity){
        index[cityLower] = zipsForCity = [];
    }
    zipsForCity.push(record);
    return index;
}, {});

const app = express(); // creates new express application

const port = process.env.PORT || 80; // use result of boolean or to assign a value to a variable, "" -> falsey
const host = process.env.HOST || '';

app.use(morgan('dev')); // accepts param for request logging format, in this case dev; called for every single request
app.use(cors()); // returns a middleware function, use cors for every request

// only called if reqMethod == GET, note order of params, opposite golang, no type checking
app.get('/hello/:name', (req, res) => { // matches calls to localhost/hello/wildcard as name (: = wildcard), uses lambda function =>
    res.send(`Hello ${req.params.name}!`); // access request params assigned in path
}); 

app.get('/zips/city/:cityName', (req, res)=> {
    let zipsForCity = zipCityIndex[req.params.cityName.toLowerCase()];
    if(!zipsForCity){
        res.status(404).send('invalid city name');
    } else {
        res.json(zipsForCity); // sends json as string to client AND sets content type header automatically (Content-Type: application-json; charset: "UTF-8")
    }
});

// note: whether or not additional functions are called dependent on whether or not previous function calls .next()

app.listen(port, host, () => {
    console.log(`server is listening at http://${host}:${port}`);
});