import { fetchFlags } from "./api.js";

const FLAG_KEYS = ["test1"];

let requestCount = 0;
let successCount = 0;
let notModifiedCount = 0;
let errorCount = 0;

export function getStats() {
  return { requestCount, successCount, notModifiedCount, errorCount };
}

async function singleRequest(id) {
  const key = FLAG_KEYS[id % FLAG_KEYS.length];
  const label = key;

  try {
    const result = await fetchFlags(key);

    if (result.notModified) {
      notModifiedCount++;
      console.log(`[304] req #${requestCount} — ${label} (not modified)`);
    } else {
      successCount++;
      const flagCount = result.data?.flags?.length ?? 0;
      console.log(`[200] req #${requestCount} — ${label} (${flagCount} flags)`);
    }
  } catch (err) {
    errorCount++;
    console.error(`[ERR] req #${requestCount} — ${label}: ${err.message}`);
  }
}

export async function generateTraffic({ intervalMs, concurrency }) {
  console.log(
    `Starting traffic: interval=${intervalMs}ms, concurrency=${concurrency}`,
  );

  const tick = () => {
    for (let i = 0; i < concurrency; i++) {
      requestCount++;
      singleRequest(requestCount);
    }
  };

  tick();
  return setInterval(tick, intervalMs);
}

export function printStats() {
  const { requestCount, successCount, notModifiedCount, errorCount } =
    getStats();
  const total = successCount + notModifiedCount + errorCount;
  console.log("\n--- Traffic Stats ---");
  console.log(`Total requests:  ${requestCount}`);
  console.log(`Successful:      ${successCount}`);
  console.log(`Not modified:    ${notModifiedCount}`);
  console.log(
    `Errors:          ${errorCount}${total > 0 ? ` (${((errorCount / total) * 100).toFixed(1)}%)` : ""}`,
  );
  console.log("--------------------\n");
}
