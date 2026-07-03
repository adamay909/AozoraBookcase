// manifest.js should be generated using the manuscript.sh script in scripts folder.

import { FILES } from "./manifest.js";

const CACHE = "azbookcase";
self.addEventListener("install", event => {
  event.waitUntil(
    Promise.all([
      caches.open(CACHE).then(cache => cache.addAll(FILES)),
      caches.open(CACHE).then(cache => {
        const req = new Request("https://cdn.jsdelivr.net/npm/idb-keyval@6/dist/umd.js");
        const options = { mode: "no-cors" };
        return fetch(req, options).then(response => {
          return cache.put(req, response.clone()).then(() => response);
        });
      })
    ]).catch(err => console.error("Install failed:", err))
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
