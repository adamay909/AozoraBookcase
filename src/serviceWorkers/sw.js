import { FILES } from "./manifest.js";

const CACHE = "azbookcase";

self.addEventListener("install", event => {
  event.waitUntil(
    caches.open(CACHE).then(cache => {
      return cache.addAll(FILES);
    })
  );
});


self.addEventListener('fetch', event => {
  event.respondWith(
    //fetch(event.request, { cache: 'no-cache' })
    fetch(event.request)
      .then(response => {
        const clone = response.clone();
        caches.open(CACHE).then(cache => cache.put(event.request, clone));
        return response;
      })
      .catch(() => {
	   console.log("offline. Serving from cache")
	   caches.match(event.request)
	  })
	  
  );
});
