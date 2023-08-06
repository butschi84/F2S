const http = require('http');
const url = require('url');

// get a random delay from 200ms to 5000ms
function getRandomDelay() {
  return Math.floor(Math.random() * (5000 - 200 + 1)) + 200;
}

const server = http.createServer((req, res) => {
  const parsedUrl = url.parse(req.url, true); // Parse the URL including query parameters
  const { pathname, query } = parsedUrl;

  const delay = query.delay ? parseInt(query.delay) : getRandomDelay();

  // non-blocking simulation
  // => container can serve other requests during delay
  switch(pathname) {
    case "/":
      res.statusCode = 200;
      res.setHeader('Content-Type', 'text/plain');
      res.end('DONE');
      break;
    case "/blocking":
       // Blocking delay using a loop
      const start = Date.now();
      while (Date.now() - start < delay) {
        // Do nothing, just wait
      }

      res.statusCode = 200;
      res.setHeader('Content-Type', 'text/plain');
      res.end('BLOCKING REQUEST DONE');
      break;
    case "/nonblocking":
      // Simulating delay before responding
      setTimeout(() => {
        res.statusCode = 200;
        res.setHeader('Content-Type', 'text/plain');
        res.end('NONBLOCKING REQUEST DONE');
      }, delay);
      break;
    case "/json":
      res.statusCode = 200;
      res.setHeader('Content-Type', 'application/json');
      res.end(req.body ? req.body : {'status': 'done'});  
      break;
  }
});

const port = 8080;
server.listen(port, () => {
  console.log(`Server running at http://0.0.0.0:${port}/`);
});
