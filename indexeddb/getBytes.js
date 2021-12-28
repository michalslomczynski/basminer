const myPromise = new Promise((resolve, reject) => {
    var DBOpenRequest = window.indexedDB.open("/idbfs", 21);
    DBOpenRequest.onerror = function (event) {
        reject(`failed to connect ${event}`);
    };

    DBOpenRequest.onsuccess = function (event) {
        db = DBOpenRequest.result;
        var objectStore = db.transaction('FILE_DATA', 'readwrite').objectStore('FILE_DATA');
        var allKeys = objectStore.getAllKeys();
        allKeys.onsuccess = function (event) {
            var prefsKey = event.target.result[1];
            var prefs = objectStore.get(prefsKey);
            prefs.onsuccess = function (event) {
                var storeBytes = event.target.result.contents;
                resolve(storeBytes);
            }
        }
    };
});