// manifest.js should be generated using the manuscript.sh script in scripts folder.

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
    fetch(event.request)
      .then(response => {
        const clone = response.clone();
        caches.open(CACHE).then(cache => cache.put(event.request, clone));
        return response;
      })
      .catch(() => 
	   caches.match(event.request)
	   .then(response => {
		if (response) {
		 return response
		}
		return new Response("Offline", {status: 404, statusText: "Not Found"})
		})
	   )
	  )
  })
