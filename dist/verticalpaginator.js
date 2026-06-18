import { paginate, setPageRanges } from "./vpcore.js";
import {
  goToPage,
  attachListeners,
  resetUI,
  setDocID,
  getCurrentPage,
} from "./vpUI.js";
import * as vpp from "./vpPersistence.js";

const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

let verticalReader, pageStrip, pageTrack, busyIndicator, loadingIndicator;
let mo;

export async function initReader(docId) {
  setConstants();
  await waitForData();
  const restored = await vpp.loadPaginatedDoc(docId);
  let current = 0;
  if (restored) {
    //found previous data
    pageStrip = document.getElementById("pageStrip");
    cleanPageStrip();
    resetUI();
    setDocID(docId);
    attachListeners();
    current = await vpp.loadSavedPage(docId);
    goToPage(current);
    setPageRanges();
    return;
  }

  pageStrip = document.getElementById("pageStrip");
  const busyIndicator = document.getElementById("busyPaginating");
  const pageTrack = document.getElementById("pageTrack");
  verticalReader.style.visibility = "visible";
  pageTrack.style.display = "none";
  busyIndicator.style.display = "inline-block";
  mo = new MutationObserver(showReader);
  mo.observe(pageStrip, { childList: true });
  await paginate();
  if (!pageStrip.isConnected) {
    //in case user has navigated away
    return;
  }
  vpp.savePaginatedDoc(docId); //no TOC at this point!
  pageTrack.style.removeProperty("display");
  busyIndicator.style.removeProperty("display");
  const currentPage = getCurrentPage();
  resetUI();
  setDocID(docId);
  attachListeners();
  goToPage(currentPage);
  requestAnimationFrame(() => {
    requestAnimationFrame(() => {});
  });
  mo.disconnect();
}

function setConstants() {
  verticalReader = document.getElementById("verticalReader");
  pageStrip = document.getElementById("pageStrip");
  pageTrack = document.getElementById("pageTrack");
  busyIndicator = document.getElementById("busyPaginating");
  loadingIndicator = document.getElementById("loading");
}

function cleanPageStrip() {
  for (let page of pageStrip.children) {
    page.style.display = "none";
  }
}
async function waitForData() {
  const images = document.querySelectorAll("img");
  const imagePromises = Array.from(images).map((img) => {
    if (img.complete && img.naturalHeight !== 0) return Promise.resolve();
    return new Promise((resolve) => {
      img.onload = resolve;
      img.onerror = resolve;
    });
  });
  await Promise.all([...imagePromises, document.fonts.ready]);
}

let uishown = false;

function showReader(mr, mo) {
  if (uishown) {
    return;
  }
  if (pageStrip.children.length > 15) {
    uishown = true;
    resetUI();
    attachListeners();
    verticalReader.style.visibility = "visible";
    goToPage(0);
    requestAnimationFrame(() => {
      requestAnimationFrame(() => {});
    });
    mo.disconnect();
  }
}
