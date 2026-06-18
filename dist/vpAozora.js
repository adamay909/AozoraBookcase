import { initReader } from "./verticalpaginator.js";
import { repaginate } from "./vpcore.js";
import {
  goToPage,
  attachListeners,
  detachListeners,
  getCurrentPage,
  resetUI,
} from "./vpUI.js";
import { savePaginatedDoc } from "./vpPersistence.js";

let pageWindow, pageStrip, navBar, currentPage, seekBar;
let handlers = false;
let docId;

const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

export async function initAozoraReader(id) {
  resetUI();
  docId = id;
  const verticalReader = document.getElementById("verticalReader");
  const loadingIndicator = document.getElementById("loading");
  verticalReader.style.visibility = "visible";
  await initReader(docId);
  initialize();
  setupTOC();
  addEventHandlers();
   requestAnimationFrame(() => {
      requestAnimationFrame(() => {});
    });
}

window.initAozoraReader = initAozoraReader;

function initialize() {
  pageWindow = document.getElementById("pageWindow");
  pageStrip = document.getElementById("pageStrip");
  navBar = document.getElementById("navBar");
  seekBar = document.getElementById("seekBar");
}

function setupTOC() {
  constructTOC();
  const toc = document.getElementById("aozoraTOC");
  if (!toc) {
    return;
  }
  const tocButton = document.querySelector("#tocLabel");
  tocButton.style.display = "inline";
  savePaginatedDoc(docId);
}

function showTOC(e) {
  e.preventDefault();
  navBar.setAttribute("inert", "inert");
  seekBar.setAttribute("inert", "inert");
  currentPage = getCurrentPage();
  detachListeners();
  const navElem = document.getElementById("aozoraTOC");
  pageStrip.style.display = "none";
  pageWindow.style.overflow = "scroll";
  navElem.style.display = "block";
}

function continueReadingAt(page) {
  const navElem = document.getElementById("aozoraTOC");
  const settingElem = document.getElementById("readerSettingsMenu");
  if (navElem) {
    navElem.style.display = "none";
  }
  settingElem.style.display = "none";
  pageWindow.style.removeProperty("overflow");
  pageStrip.style.removeProperty("display");
  navBar.removeAttribute("inert", "inert");
  seekBar.removeAttribute("inert", "inert");
  attachListeners();
  goToPage(page);
}

function jumpToTarget(e) {
  e.preventDefault();
  const link = e.target.closest("a.toclink");
  if (!link) {
    return;
  }
  const url = new URL(link.href);
  if (!url) {
    return;
  }
  const destination = document.querySelector(url.hash).closest(".page");
  const idx = Number(destination.id.slice(4));
  continueReadingAt(idx);
}

function backToReading(e) {
  e.preventDefault();
  continueReadingAt(currentPage);
}

async function tocEventListeners(e) {
  let target = e.target.closest("a");

  if (target) {
    if (target.classList.contains("toclink")) {
      jumpToTarget(e);
      return;
    }
  }

  target = e.target.closest("button");

  if (target) {
    switch (target.id) {
      case "backToReading":
        backToReading(e);
        return;

      case "openTOC":
        showTOC(e);
        break;

      case "readerSettings":
        showSettings(e);
        break;

      case "cancel":
        pageWindow.style.cssText = currentSettings;
        backToReading(e);
        break;

      case "doit":
        if (currentSettings != newSettings) {
          pageWindow.style.cssText = newSettings;
          const cover = document.getElementById("loading");
          cover.style.display = "flex";
          currentPage = await repaginate(currentPage);
          savePaginatedDoc(docId);
          sleep(0);
          requestAnimationFrame(() => {
            requestAnimationFrame(() => {});
          });
          cover.style.display = "none";
          resetUI();
          backToReading(e);
        }
        backToReading(e);
        break;
      case "showDebugInfo":
        const elem = document.getElementById("debugInfo");
        if (elem.innerText.trim().length > 0) {
          elem.replaceChildren();
          pageWindow.classList.remove("debug");
          return;
        }
        const css = window.getComputedStyle(pageWindow);
        pageWindow.classList.add("debug");
        const w = css.width;
        const h = css.height;
        const fs = css.fontSize;
        const overFlowPages = document.querySelectorAll(".overflow");
        let text =
          `width: ${w}` +
          "<br>" +
          `height: ${h}` +
          "<br>" +
          `font-size: ${fs}` +
          "<br>";
        if (overFlowPages.length > 0) {
          let t2 = "";
          for (let k of overFlowPages) {
            t2 = t2 + k.id + ", ";
          }
          text = text + "overflow: " + t2 + "<br>";
        }
        elem.innerHTML = text;
        break;
    }
  }
}

function constructTOC() {
  const tocFrag = document.querySelector("#documentTOC");
  if (!tocFrag) {
    return;
  }
  const tocContent = tocFrag.content;

  //count the entries in TOC
  const tocWalker = document.createTreeWalker(
    tocContent,
    NodeFilter.SHOW_ELEMENT,
  );

  let count = 0;

  for (
    let elem = tocWalker.nextNode();
    elem != null;
    elem = tocWalker.nextNode()
  ) {
    if (elem.tagName === "LI") {
      count++;
    }
    if (elem.tagName === "A") {
      elem.classList.add("toclink"); //add class in case we show TOC
    }
  }

  //no point showing a TOC that is too short
  if (count < 2) {
    const tmp = document.createElement("DIV");
    tmp.appendChild(tocContent);
    tmp.remove();
    return;
  }

  //the first listed item is the title of the book. remove link
  //for that
  tocWalker.currentNode = tocWalker.root;
  for (
    let elem = tocWalker.nextNode();
    elem != null;
    elem = tocWalker.nextNode()
  ) {
    if (elem.tagName === "A") {
      const title = elem.innerText;
      const item = elem.parentElement;
      const titleElem = document.createElement("H3");
      titleElem.innerText = title;
      const tocHeader = document.createElement("H4");
      tocHeader.innerText = "目次";
      titleElem.style.marginBottom = "1em";
      tocHeader.style.marginBottom = "2em";
      item.replaceChild(titleElem, elem);
      titleElem.after(tocHeader);
      break;
    }
  }
  const navElem = document.createElement("DIV");
  navElem.appendChild(tocContent);

  navElem.id = "aozoraTOC";
  navElem.style.display = "none";
  navElem.writingMode = "horizontal-tb";
  navElem.overflow = "scroll";

  const buttonBar = document.createElement("DIV");
  buttonBar.style.display = "flex";
  buttonBar.style.position = "sticky";
  buttonBar.style.top = "10px";

  const backButton = document.createElement("BUTTON");
  backButton.appendChild(document.createTextNode("読書に戻る"));
  backButton.style.marginLeft = "auto";
  backButton.style.marginRight = "10px";
  backButton.style.marginBottom = "10px";
  backButton.id = "backToReading";

  buttonBar.prepend(backButton);

  navElem.prepend(buttonBar);

  pageWindow.appendChild(navElem);
}

function addEventHandlers() {
  if (!handlers) {
    document.addEventListener("click", tocEventListeners, { passive: false });
    document.addEventListener("change", readSettings, { passive: false });

    handlers = true;
  }
}

let currentSettings, newSettings;

function showSettings(e) {
  e.preventDefault();
  navBar.setAttribute("inert", "inert");
  seekBar.setAttribute("inert", "inert");
  currentPage = getCurrentPage();
  detachListeners();
  currentSettings = pageWindow.style.cssText;
  newSettings = currentSettings;
  const fontSizes = document.getElementById("fontSizes");
  switch (pageWindow.style.fontSize) {
    case "0.85rem":
      fontSizes.value = "0";
      break;
    case "1.2rem":
      fontSizes.value = "2";
      break;
    default:
      fontSizes.value = 1;
  }
  const lineHeight = document.getElementById("lineHeight");
  switch (pageWindow.style.lineHeight) {
    case "1.75":
      lineHeight.value = "0";
      break;
    case "2.2":
      lineHeight.value = "2";
      break;
    default:
      lineHeight.value = "1";
  }

  const overFlowPages = document.querySelectorAll(".overflow");
  if (overFlowPages.length > 0) {
    document.getElementById("showDebugInfo").style.display = "block";
  }
  const settingElem = document.getElementById("readerSettingsMenu");
  settingElem.style.display = "block";
}

function readSettings(e) {
  const target = e.target.closest("input");

  if (!target) {
    return;
  }
  switch (target.id) {
    case "fontSizes":
      switch (target.value) {
        case "0":
          pageWindow.style.fontSize = "0.85rem";
          break;
        case "1":
          pageWindow.style.fontSize = "1rem";
          break;
        case "2":
          pageWindow.style.fontSize = "1.2rem";
          break;
      }
      break;

    case "lineHeight":
      switch (target.value) {
        case "0":
          pageWindow.style.lineHeight = "1.75";
          break;
        case "1":
          pageWindow.style.lineHeight = "1.9";
          break;
        case "2":
          pageWindow.style.lineHeight = "2.2";
          break;
      }
  }
  newSettings = pageWindow.style.cssText;
}
