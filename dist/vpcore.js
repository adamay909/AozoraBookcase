let pageWidth, columnHeight, maxWidth;
let pageWindow, pageStrip, sandbox;
let pages = [];
let extracss = document.createElement("style");
let paginated = false;
const TRANSPARENT_CONTAINERS = new Set([
  "SECTION",
  "ARTICLE",
  "ASIDE",
  "MAIN",
  "HEADER",
  "FOOTER",
]);

const extraCSS = `
sup {
vertical-align: baseline;
}
sub {
vertical-align: baseline;
}
`;

const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

export async function paginate() {
  waitForData();
  console.log("BEGIN PAGINATION");
  initialize();
  const src = document.createElement("DIV");
  src.appendChild(document.getElementById("documentContent").content);
  toBlocks(src);
  await getPageRanges();
  cleanup();
  console.log("END PAGINATION");
}

//currentPage is the page opened now. Return value is approximate
//page number after repagination.
export async function repaginate(currentPage = 0) {
  if (!paginationComplete()) {
    await paginationFinished();
  }
  let startRange = pages[currentPage]; //this is where we are
  let startMark = document.createRange();
  startMark.setStart(startRange.startContainer, startRange.startOffset);
  startMark.setEnd(startRange.endContainer, startRange.endOffset);
  console.log("START RE-PAGINATION");
  initialize();
  sandbox = document.getElementById("paginationSandbox");
  await getPageRanges();
  let i = 0;
  for (i = 0; i < pages.length - 1; i++) {
    if (startMark.compareBoundaryPoints(Range.START_TO_START, pages[i]) < 1) {
      break;
    }
  }
  currentPage = i - 1;
  cleanup();
  console.log("END RE-PAGINATION");
  return currentPage;
}

async function paginationFinished() {
  while (!paginationComplete()) {
    await sleep(5);
  }
  return;
}

function paginationComplete() {
  return pages[pages.length - 1] == "finished";
}

function initialize() {
  pages = [];
  pageWindow = document.getElementById("pageWindow");
  pageStrip = document.getElementById("pageStrip");
  pageStrip.replaceChildren();
  pageStrip.classList.add("page-strip");
  pageWindow.style.removeProperty("width");
  pageWidth = pageWindow.clientWidth;
  columnHeight = pageWindow.clientHeight;
  sandbox = document.getElementById("paginationSandbox");
  if (!sandbox) {
    sandbox = document.createElement("DIV");
    sandbox.id = "paginationSandbox";
    pageWindow.appendChild(sandbox);
  }
  sandbox.classList.add("vpSandbox");
  const page = document.createElement("DIV");
  pageStrip.appendChild(page);
  const pageBody = document.createElement("DIV");
  page.appendChild(pageBody);
  const pageBodyCss = window.getComputedStyle(pageBody);
  page.remove();
  sandbox.classList.add("vpShow");
  if (pageBodyCss.cssText.trim()) {
    sandbox.style.cssText = pageBodyCss.cssText; //reflect all styles of page body.
  }
  sandbox.classList.remove("vpShow");
  extracss.textContent = extraCSS;
  extracss.id = "vpcss";
  document.head.appendChild(extracss);
}

//src is the root node of the document we want to paginate. The result
//will be in sandbox.
function toBlocks(src) {
  src.classList.add("vpSrc");
  sandbox.replaceChildren();
  const walker = document.createTreeWalker(src, NodeFilter.SHOW_ALL, filter);
  function filter(node) {
    for (let e = node.parentNode; e != null; e = e.parentNode) {
      if (e.nodeType == Node.ELEMENT_NODE && e.tagName == "FIGURE") {
        //avoid double counting img inside FIGURE
        return NodeFilter.FILTER_REJECT;
      }
    }
    return NodeFilter.FILTER_ACCEPT;
  }

  let newBlock = true; //whether next node starts new block
  let startContainer, startOffset, caretNode, tmprange;
  while (walker.nextNode()) {
    caretNode = walker.currentNode;
    if (newBlock) {
      startContainer = caretNode;
      startOffset = 0;
      newBlock = false;
    }

    switch (caretNode.tagName) {
      case "DIV":
        if (caretNode.classList.contains("pageBreak")) {
          //in this case, we just end the previous block and create a new one
          closePreviousBlock();
          newBlock = true;
        }
        break;

      case "FIGURE": //create block that contains just the figure
        closePreviousBlock();
        tmprange = document.createRange();
        tmprange.selectNode(caretNode);
        addNewBlock(tmprange);
        newBlock = true;
        break;

      case "IMG": //ditto
        if (!caretNode.classList.contains("illustration")) {
          break;
        }
        if (caretNode.closest("FIGURE") != null) {
          break;
        }
        closePreviousBlock();
        tmprange = document.createRange();
        tmprange.selectNode(caretNode);
        addNewBlock(tmprange);
        newBlock = true;
        break;
    }
  }
      closePreviousBlock();

  //HOISTED FUNCTIONS
  function closePreviousBlock() {
    let endContainer = walker.previousNode();
    walker.currentNode = caretNode;
    let endOffset =
      endContainer.nodeType === Node.TEXT_NODE
        ? endContainer.textContent.length
        : endContainer.childNodes.length;
    const blockRange = makeRange(
      startContainer,
      startOffset,
      endContainer,
      endOffset,
    );
    if (!blockRange.toString().trim()) {
      return; //reject empty pages
    }
    addNewBlock(blockRange);
  }

  function addNewBlock(range) {
    const tree = getTree(range, ".vpSrc", false);
    const block = document.createElement("DIV");
    block.classList.add("vpBlock");
    block.replaceChildren(...tree.childNodes);
    sandbox.appendChild(block);
  }
}

export async function setPageRanges() {
  pages = [];
  pageWindow = document.getElementById("pageWindow");
  pageStrip = document.getElementById("pageStrip");
  sandbox = document.getElementById("paginationSandbox");
  pageWidth = pageWindow.clientWidth;
  columnHeight = pageWindow.clientHeight;
  getPageRanges(false);
}

async function getPageRanges(render = true) {
  const pageRanges = [];
  let pageNum = 0;
  let walker, caretNode;
  let startNode, startOffset, endNode, endOffset;
  sandbox.classList.add("vpShow");
  //display only the block under consideration
  for (let i = 0; i < sandbox.children.length; i++) {
    layoutBlock(i);
    const paginationRoot = sandbox.children[i];
    //first check if are dealing with known one-page blocks
    if (blockIsSinglePage(paginationRoot)) {
      pages.push(getSinglePage(paginationRoot));
      pageNum++;
      if (render) {
        renderPage(pageNum - 2);
      } //we need at least two pages to render
      await sleep(0);
      continue;
    }
    //we are dealing with a potentially multi-page block
    //let's get the boundaries for each page. We simply look
    //maximal rects that fit inside a page width.
    let xMax = 0; //xMax 'grows' in negative direction

    walker = document.createTreeWalker(
      paginationRoot,
      NodeFilter.SHOW_ALL,
      paginationFilter,
    );

    startNode = walker.nextNode();
    startOffset = 0;
    walker.currentNode = paginationRoot;
    let prange = document.createRange();
    prange.setStart(startNode, startOffset);

    while ((caretNode = walker.nextNode())) {
      const testrange = document.createRange(caretNode);
      testrange.selectNode(caretNode);
      const rectangles = testrange.getClientRects();
      for (let i = 0; i < rectangles.length; i++) {
        if (rectangles[i].left > xMax) {
          continue;
        }

        endNode = caretNode;
        endOffset = findEndOffset(endNode, xMax);
        if (endOffset == 0) {
          endNode = walker.previousNode();
          endOffset = endNode.length;
          startNode = walker.nextNode(); //this also brings walker back
          startOffset = 0;
        } else {
          startNode = caretNode;
          startOffset = endOffset;
        }
        prange.setEnd(endNode, endOffset);
        pages.push(prange);
        prange = document.createRange();
        prange.setStart(startNode, startOffset);
        xMax = rectangles[i].right - pageWidth;
        pageNum++;
        if (render) {
          renderPage(pageNum - 2); //we need at least two pages
          await sleep(0);
        }
      }
    }
    //don't forget last page of block
    endNode = walker.previousNode();
    endOffset = endNode.length;
    prange.setEnd(endNode, endOffset);
    pages.push(prange);
    pageNum++;
    if (render) {
      renderPage(pageNum - 2);
      await sleep(0);
    }
  }
  if (render) {
    renderPage(pageNum - 1); //render final page
    await sleep(0);
  }
  //clean up before returning
  sandbox.classList.remove("vpShow");
  for (let elem of sandbox.querySelectorAll(".vpShow")) {
    elem.classList.remove("vpShow");
  }
  pages.push("finished");
  return;

  //////END OF MAIN PART OF FUNCTION
  //the rest are hoisted functions:
  function paginationFilter(node) {
    if (node.nodeType === Node.ELEMENT_NEDE) {
      if (node.tagName === "RT") {
        return NodeFilter.FILTER_REJECT;
      }
    }
    if (node.nodeType === Node.TEXT_NODE) {
      return NodeFilter.FILTER_ACCEPT;
    }

    if (node.nodeType == Node.ELEMENT_NODE && node.tagName == "HR") {
      return NodeFilter.FILTER_ACCEPT;
    }

    return NodeFilter.FILTER_SKIP;
  }

  function layoutBlock(num) {
    if (num > 0) {
      sandbox.children[num - 1].classList.remove("vpShow");
    }
    sandbox.children[num].classList.add("vpShow");
    //wait for reflow
     sleep(1);
    requestAnimationFrame(() => {
      requestAnimationFrame(() => {});
    });
  }

  function blockIsSinglePage(paginationRoot) {
    if (paginationRoot.querySelector("FIGURE") != null) {
      return true;
    }
    if (paginationRoot.querySelector("IMG.illustration") != null) {
      return true;
    }
    return false;
  }

  function getSinglePage(paginationRoot) {
    const tmprange = document.createRange();
    tmprange.selectNode(paginationRoot.children[0]);
    return tmprange;
  }

  function findEndOffset(endNode, xMax) {
    //find offset of endNode
    let low = 0;
    let high = endNode.length;
    const tempRange = document.createRange();
    while (low < high) {
      let mid = Math.floor((low + high) / 2);
      tempRange.setStart(endNode, mid);
      tempRange.setEnd(endNode, mid + 1);
      const rects = tempRange.getClientRects();
      const rect = tempRange.getBoundingClientRect();
      if (rect.left < xMax) {
        high = mid;
      } else {
        low = mid + 1;
      }
    }
    return low;
  }
}

function renderPage(i) {
  if (i < 0) {
    return;
  }
  if (i > 0) {
    markSplit2(pages[i]);
  }
  if (i < pages.length - 1) {
    markSplit1(pages[i], pages[i + 1]);
  }
  commitPage(pages[i]);
}

function markSplit2(pageRange) {
  const startC = closestContainer(pageRange.startContainer);
  if (startC.classList.contains("vpSplit1")) {
    startC.classList.remove("vpSplit1");
    startC.classList.add("vpSplit2");
  }
}

function markSplit1(pageRange1, pageRange2) {
  const container1 = closestContainer(pageRange1.endContainer);
  const container2 = closestContainer(pageRange2.startContainer);
  if (container1 != container2) {
    return;
  }
  const testpage1 = getTree(pageRange1, ".vpBlock", false);
  if (lastElemType(testpage1) === "BR") {
    return;
  }
  container1.classList.add("vpSplit1");
}

//returns closest ancestor that is a <p>, <div>, <ol>, <ul>
//these are container elements that are likely to need adjustments
//when broken into two during pagination
function closestContainer(node) {
  for (node = node; node.parentNode !== null; node = node.parentNode) {
    if (
      node.nodeType === Node.ELEMENT_NODE &&
      (node.tagName == "P" ||
        node.tagName == "DIV" ||
        node.tagName == "OL" ||
        node.tagName == "UL")
    ) {
      return node;
    }
  }
  return null;
}

function commitPage(pageRange) {
  //each page must have a partial tree from the root down that covers the
  //nodes present in the range.
  const tree = getTree(pageRange, ".vpBlock", true);
  tree.removeAttribute("style");
  //for (const elem of tree.getElementsByTagName("*")) {
  // elem.style.removeProperty("display");
  //}
  const page = makePage();
  page.appendChild(tree);
  page.style.display = "none";
  tree.classList.remove("vpBlock");
  tree.classList.add("page-body");
  pageStrip.appendChild(page);
  page.id = "page" + String(pageStrip.children.length - 1).padStart(2, "0");
  checkOverflow(page);
}

// create one blank page. Not attached to anything.
function makePage() {
  const page = document.createElement("div");
  page.className = "page";
  page.style.width = pageWidth + "px";
  page.style.height = columnHeight + "px";
  return page;
}

function makeRange(startContainer, startOffset, endContainer, endOffset) {
  const newRange = document.createRange();
  newRange.setStart(startContainer, startOffset);
  newRange.setEnd(endContainer, endOffset);
  return newRange;
}

//This returns a partial tree rooted in an element that matches the selector. You must make
//sure that something up the tree of the range matches the selector! If simplify is set,
//tree is shortened to skip containers that match TRANSPARENT_CONTAINERS
function getTree(range, selector, simplify) {
  let tip = range.commonAncestorContainer;
  while (tip.nodeType != Node.ELEMENT_NODE) {
    tip = tip.parentNode;
  }
  let tree = tip.cloneNode(false);
  tree.style.removeProperty("display");
  const frag = range.cloneContents();
  tree.appendChild(frag); //this gives us the partial tree down from commonAncestorContainer (or its parent)
  let newRoot;
  while (!tip.matches(selector)) {
    tip = tip.parentNode;
    newRoot = tip.cloneNode(false);
    newRoot.style.removeProperty("display");
    newRoot.appendChild(tree);
    tree = newRoot;
  }
  if (simplify) {
    tree = cleanTree(tree);
  }
  return tree;
}

//return a shortened tree rooted at node by removing nodes matching TRANSPARENT_CONTAINERS
function cleanTree(node) {
  const walker = document.createTreeWalker(
    node,
    NodeFilter.SHOW_ALL,
    cleanFilter,
  );

  const root = node.cloneNode(false);
  const nodeMap = new Map();
  nodeMap.set(node, root);
  let appendPoint = root;
  let counter = 0;
  while ((node = walker.nextNode())) {
    const clone = node.cloneNode(false);
    appendPoint = nodeMap.get(walker.parentNode());
    walker.currentNode = node;
    appendPoint.appendChild(clone);
    nodeMap.set(walker.currentNode, clone);
  }
  return root;

  function cleanFilter(node) {
    if (
      node.nodeType == Node.ELEMENT_NODE &&
      TRANSPARENT_CONTAINERS.has(node.tagName)
    ) {
      return NodeFilter.FILTER_SKIP;
    }
    return NodeFilter.FILTER_ACCEPT;
  }
}

function waitForData() {
  const images = document.querySelectorAll("img");
  const imagePromises = Array.from(images).map((img) => {
    if (img.complete && img.naturalHeight !== 0) return Promise.resolve();
    return new Promise((resolve) => {
      img.onload = resolve;
      img.onerror = resolve;
    });
  });
  Promise.all([...imagePromises, document.fonts.ready]);
}

function checkOverflow(page) {
  page.style.visibility = "hidden";
  page.style.removeProperty("display");
  sleep(0);
  const w = page.scrollWidth;
  page.style.display = "none";
  page.style.removeProperty("visibility");
  if (w <= pageWidth + 10) {
    return;
  }
  page.style.width = w + "px";
  page.classList.add("overflow");
  if (w > maxWidth) {
    maxWidth = w;
  }
}

export function cleanup() {
  extracss.remove();
}

function lastElemType(tree) {
  const walker = document.createTreeWalker(tree, NodeFilter.SHOW_ALL);
  while (walker.nextNode()) {}
  const target = walker.previousNode();
  if (target.nodeType !== Node.ELEMENT_NODE) {
    return "";
  }
  return target.tagName;
}

function isLastElem(elem, container) {
  const walker = document.createTreeWalker(container, NodeFilter.SHOW_ALL);
  while (walker.nextNode()) {}
  const target = walker.currentNode;
  return target === elem;
}

function isFirstElem(elem, container) {
  const walker = document.createTreeWalker(container, NodeFilter.SHOW_ALL);
  while (walker.firstChild()) {}
  const target = walker.currentNode;
  return target === elem;
}
