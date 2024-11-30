self.addEventListener('install', (event) => {
    self.skipWaiting();
});

self.addEventListener('activate', (event) => {
    event.waitUntil(clients.claim());
});

// No caching implementation as requested
self.addEventListener('fetch', (event) => {
    event.respondWith(fetch(event.request));
});
