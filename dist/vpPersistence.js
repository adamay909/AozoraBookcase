let pageWindow, pageStrip;

pageWindow = document.getElementById("pageWindow");
pageStrip = document.getElementById("pageStrip");

export function getDocId() {
  const tpl = document.getElementById("documentContent");
  const metadata = tpl.content.querySelector("[data-docid]");
  return metadata ? metadata.dataset.docid : null;
}

//you need idb-keyval installed globally:
//<script src="https://cdn.jsdelivr.net/npm/idb-keyval@6/dist/umd.js"></script>

function idbGet(key) {
  return idbKeyval.get(key);
}

function idbSet(key, value) {
  return idbKeyval.set(key, value);
}

// Save current page number.
export async function saveState(id, current) {
  const docStateKey = id + "page";
  try {
    await idbSet(docStateKey, current);
  } catch (_) {}
}

// Load the saved page number, check consistency with current pagination.
export async function loadSavedPage(id) {
  pageStrip = document.getElementById("pageStrip");
  console.log("retrieving position");
  try {
    const docStateKey = id + "page";
    const pagenum = await idbGet(docStateKey);
    const idx = parseInt(pagenum, 10);
    console.log("found", idx);
    return Number.isFinite(idx) && idx >= 0 && idx < pageStrip.children.length
      ? idx
      : 0;
  } catch (_) {
    return 0;
  }
}

// Save paginated HTML with page dimensions
export async function savePaginatedDoc(id) {
  console.log("attempt to save", id);
  pageWindow = document.getElementById("pageWindow");
  pageStrip = document.getElementById("pageStrip");
  try {
    await idbSet(id, {
      html: pageWindow.outerHTML,
    });
    console.log("saved file");
  } catch (_) {}
}

// Restore a previously paginated document.
// Returns true if a valid, dimension-matched cache was found and loaded;
// returns false if pagination must be run from scratch.
export async function loadPaginatedDoc(id) {
  pageWindow = document.getElementById("pageWindow");
  try {
    const entry = await idbGet(id);
    if (!entry || !entry.html) return false;

    console.log("found paginated data");
    pageWindow.outerHTML = entry.html;
    return true;
  } catch (_) {
    return false;
  }
}
