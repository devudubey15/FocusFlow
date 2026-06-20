// FocusFlow Background Script
const hostName = "com.devu.focusflow";
let port = null;

function connect() {
    port = chrome.runtime.connectNative(hostName);
    port.onMessage.addListener((msg) => {
        console.log("Received from native host:", msg);
    });
    port.onDisconnect.addListener(() => {
        console.log("Disconnected from native host");
        port = null;
    });
}

function sendUrl(url) {
    if (!port) connect();
    try {
        port.postMessage({ url: url, timestamp: new Date().toISOString() });
    } catch (e) {
        console.error("Error sending message to native host:", e);
    }
}

// Track tab updates
chrome.tabs.onUpdated.addListener((tabId, changeInfo, tab) => {
    if (changeInfo.status === 'complete' && tab.active) {
        sendUrl(tab.url);
    }
});

// Track tab switches
chrome.tabs.onActivated.addListener(async (activeInfo) => {
    const tab = await chrome.tabs.get(activeInfo.tabId);
    if (tab.url) {
        sendUrl(tab.url);
    }
});
