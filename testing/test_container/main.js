const http = require('http');

// get a random delay from 200ms to 5000ms
function getRandomDelay() {
  return Math.floor(Math.random() * (5000 - 200 + 1)) + 200;
}

const server = http.createServer((req, res) => {
  // non-blocking simulation
  // => container can serve other requests during delay
  if (req.url === '/') {
    const delay = getRandomDelay();

    // Simulating delay before responding
    setTimeout(() => {
      res.statusCode = 200;
      res.setHeader('Content-Type', 'text/plain');
      res.end('NONBLOCKING REQUEST DONE');
    }, delay);

  // blocking simulation
  // => container CAN NOT serve other requests during delay
  }else if(req.url === '/blocking') {
    const delay = getRandomDelay();

    // Blocking delay using a loop
    const start = Date.now();
    while (Date.now() - start < delay) {
      // Do nothing, just wait
    }

    res.statusCode = 200;
    res.setHeader('Content-Type', 'text/plain');
    res.end('BLOCKING REQUEST DONE');
  } else {
    res.statusCode = 404;
    res.end('Not found');
  }
});

const port = 3000;
server.listen(port, () => {
  console.log(`Server running at http://localhost:${port}/`);
});
