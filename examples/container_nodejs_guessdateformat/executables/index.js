const http = require('http');
const url = require('url');
const guessFormat = require('moment-guess');

const server = http.createServer((req, res) => {
  const parsedUrl = url.parse(req.url, true);
  const { pathname, query } = parsedUrl;

  try {
    switch(pathname) {
      case "/":
        if (req.method === 'POST') {
          let requestBody = '';
          req.on('data', (chunk) => {
            requestBody += chunk.toString();
          });
      
          req.on('end', () => {
            try {
              res.statusCode = 200;
              res.setHeader('Content-Type', 'application/json'); // Set JSON content type
              res.end(JSON.stringify({
                result: guessFormat(requestBody, "strftime")
              })); // Send JSON response with the same data
            } catch (error) {
              res.statusCode = 400; // Bad Request if JSON parsing fails
              res.setHeader('Content-Type', 'application/json');
              res.end(JSON.stringify({ error: 'Invalid JSON' }));
            }
          });
        }else{
          res.statusCode = 200;
          res.setHeader('Content-Type', 'application/json');
          res.end(req.body ? req.body : JSON.stringify({'result': ''})); 
        }
        break;
    }
  }catch(ex){
  }
});

const port = 8080;
server.listen(port, () => {
  console.log(`Server running at http://0.0.0.0:${port}/`);
});
