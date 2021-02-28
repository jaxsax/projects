// https://github.com/mozilla/readability#nodejs-usage

var fs = require('fs');
var { Readability } = require('@mozilla/readability');
var { JSDOM } = require('jsdom');

let buffer = fs.readFileSync(process.argv[2])

var doc = new JSDOM(buffer, {
  url: "https://www.example.com/the-page-i-got-the-source-from"
});
let reader = new Readability(doc.window.document);
let article = reader.parse();

console.log(article.textContent)
