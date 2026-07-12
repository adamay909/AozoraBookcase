
async function loadFileFromIDB(name) {
    try {
      const file = await idbKeyval.get(name);
      if (!file || !file.data) {
	   return "";
	  }
	 return file.data;
    } catch (error) {
	  console.log("couldn't load", name);
	  console.error(error);
      return "";
    }
  }

async function saveFileToIDB(name, content) {
    try {
      await idbKeyval.set(name, {
	    data: content,
      });
    } catch (error) {
	console.log("couldn't save", name);
	  console.error(error);
	}
 }

async function deleteFile(name) {
 try {
   	await idbKeyval.del(name)
 } catch {
  console.log("file", name, "does not exist")
 }
}

async function fileExists(name) {
   const keys = await idbKeyval.keys()
   return keys.includes(name)
}

async function clearIDB() {
 await idbKeyval.clear()
}


//given event e, return closest element that meets selector
function eventTargetID(e, selector) {
const elem = e.target.closest(selector)
 if (!elem) {
  return ""
 }
 return elem.id
}

function eventTargetElem(e, selector) {
const elem = e.target.closest(selector)
 if (!elem) {
  return null
 }
 return elem
}

async function clearData() {

    localStorage.clear();//we don't use localStorage but just in case.

    sessionStorage.clear();//ditto
  
    await clearIDB();//clear indexDB

    const CACHE = "azbookcase"; //clear service worker cache
    await  caches.delete(CACHE);

    const registrations = await navigator.serviceWorker.getRegistrations();
    await Promise.all(registrations.map(r => r.unregister())); //unregister service workers

}


window.loadFileFromIDB = loadFileFromIDB
window.saveFileToIDB = saveFileToIDB
window.fileExists = fileExists
window.clearIDB = clearIDB
window.eventTargetID = eventTargetID
window.eventTargetElem = eventTargetElem
window.deleteFile = deleteFile
window.clearData = clearData
