import { config } from '../env.js';
import { generateTraffic, printStats } from './traffic.js';

async function main() {
  console.log('M2M Traffic Generator');
  console.log(`Target: ${config.API_BASE_URL}`);
  console.log(`Auth0 domain: ${config.AUTH0_DOMAIN}`);
  console.log(`Audience: ${config.AUTH0_AUDIENCE}`);
  console.log(`Scope: ${config.AUTH0_SCOPE}`);
  console.log('');

  const intervalId = await generateTraffic({
    intervalMs: config.TRAFFIC_INTERVAL_MS,
    concurrency: config.TRAFFIC_CONCURRENCY,
  });

  let statsInterval = setInterval(printStats, 30000);

  process.on('SIGINT', () => {
    console.log('\nShutting down...');
    clearInterval(intervalId);
    clearInterval(statsInterval);
    printStats();
    process.exit(0);
  });

  process.on('SIGTERM', () => {
    console.log('\nShutting down...');
    clearInterval(intervalId);
    clearInterval(statsInterval);
    printStats();
    process.exit(0);
  });
}

main().catch((err) => {
  console.error('Fatal error:', err);
  process.exit(1);
});
