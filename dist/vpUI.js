import { saveState } from "./vpPersistence.js";

let current = 0;
let isAnimating = false;
let docId = null;
let pagesDisplayed = -1;
let handlers = null;

let swipeTouchId = null;
let swipeStartX = 0;
const SWIPE_THRESHOLD = 30;

let pageWindow,
  pageStrip,
  loading,
  seekBar,
  pageTrack,
  readingPosition,
  seekLabel,
  verticalReader;

const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

export function getCurrentPage() {
  return current;
}

export function resetUI() {
  pageWindow = document.getElementById("pageWindow");
  pageStrip = document.getElementById("pageStrip");
  loading = document.getElementById("loading");
  seekBar = document.getElementById("seekBar");
  readingPosition = document.getElementById("readingPosition");
  seekLabel = document.getElementById("seekLabel");
  verticalReader = document.getElementById("verticalReader");
  current = 0;
  isAnimating = false;
  docId = null;
  pagesDisplayed = -1;
  swipeTouchId = null;
  swipeStartX = 0;
}

function updateUI(idx) {
  const total = pageStrip.children.length;
  if (total === 0) {
    return;
  }
  current = idx;
  readingPosition.min = 0;
  readingPosition.max = total - 1;
  readingPosition.value = current;
  const page = idx + 1;
  seekLabel.textContent = page + "/" + total;
  saveState(docId, current);
}

function pageElemOf(idx) {
  const pageID = "#page" + String(idx).padStart(2, "0");
  return pageWindow.querySelector(pageID);
}

export function goToPage(idx) {
  if (idx < 0 || idx >= pageStrip.children.length) {
    return;
  }
  const oldpage = pageElemOf(current);
  current = idx;
  const newpage = pageElemOf(idx);
  updateUI(current);
  oldpage.style.display = "none";
  newpage.style.removeProperty("display");
  pageElemOf(idx).scrollIntoView({
    behavior: "instant",
    block: "nearest",
  });
}

function move(dir) {
  if (isAnimating) return;
  const next = current + dir;
  if (next < 0 || next >= pageStrip.children.length) return;
  isAnimating = true;
  goToPage(next);
  isAnimating = false;
  return;
}

// Seek helpers (shared between attach and detach)

// Build the handlers object
// Called once; returns a plain object whose values are the named functions
// that will be passed to both addEventListener and removeEventListener.

let timeOut = performance.now();

function buildHandlers() {
  return {
    // Keyboard
    keydown(e) {
      switch (e.key) {
        case "ArrowLeft":
        case "ArrowDown":
        case "PageDown":
          e.preventDefault();
          move(1);
          break;
        case "ArrowRight":
        case "ArrowUp":
        case "PageUp":
          e.preventDefault();
          move(-1);
          break;
        case " ":
          e.preventDefault();
          move(e.shiftKey ? -1 : 1);
          break;
        case "Home":
          e.preventDefault();
          goToPage(0);
          break;
        case "End":
          e.preventDefault();
          goToPage(pageStrip.children.length - 1);
          break;
      }
    },

    //Mouse wheel. Touchpad scrolls are also reported as wheel events.
    mouseWheel(e) {
      e.preventDefault();
      if (performance.now() - timeOut < 250) {
        return;
      }
      timeOut = performance.now();
      if (Math.abs(e.deltaY / e.deltaX) > 3) {
        if (e.deltaY > 0) {
          // wheel down forward
          move(1);
        }
        if (e.deltaY < 0) {
          // wheel up backward
          move(-1);
        }
        return;
      }

      if (e.deltaX > 0) {
        // wheel down forward
        move(-1);
      }
      if (e.deltaX < 0) {
        // wheel up backward
        move(1);
      }
    },

    // Touch swipe
    touchstart(e) {
      if (swipeTouchId !== null) return;
      const t = e.changedTouches[0];
      swipeTouchId = t.identifier;
      swipeStartX = t.clientX;
    },

    touchend(e) {
      let t = null;
      for (const touch of e.changedTouches) {
        if (touch.identifier === swipeTouchId) {
          t = touch;
          break;
        }
      }
      if (!t) return;
      swipeTouchId = null;
      const dx = t.clientX - swipeStartX;
      if (Math.abs(dx) >= SWIPE_THRESHOLD) {
        move(dx > 0 ? 1 : -1);
      }
    },

    touchcancel() {
      swipeTouchId = null;
    },

    seekPage(e) {
      const idx = e.target.value;
      goToPage(Number(idx));
    },
  };
}

export function attachListeners() {
  if (handlers) return; // already attached
  handlers = buildHandlers();

  document.addEventListener("keydown", handlers.keydown);
  document.addEventListener("wheel", handlers.mouseWheel, { passive: false });
  document.addEventListener("touchstart", handlers.touchstart, {
    passive: true,
  });
  document.addEventListener("touchend", handlers.touchend, {
    passive: true,
  });
  document.addEventListener("touchcancel", handlers.touchcancel, {
    passive: true,
  });
  readingPosition.addEventListener("input", handlers.seekPage, {
    passive: true,
  });
}

export function detachListeners() {
  if (!handlers) return; // nothing to remove

  document.removeEventListener("keydown", handlers.keydown);
  document.removeEventListener("wheel", handlers.mouseWheel);
  document.removeEventListener("touchstart", handlers.touchstart);
  document.removeEventListener("touchend", handlers.touchend);
  document.removeEventListener("touchcancel", handlers.touchcancel);
  readingPosition.removeEventListener("input", handlers.seekPage);

  handlers = null; // allow re-attach later
}


export function setDocID(id) {
  docId = id;
}
window.attachReaderListeners = attachListeners;
window.detachReaderListeners = detachListeners;
