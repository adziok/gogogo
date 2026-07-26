import dotenv from "dotenv";
dotenv.config();

function required(name) {
  const value = process.env[name];
  if (!value) {
    console.error(`Missing required environment variable: ${name}`);
    process.exit(1);
  }
  return value;
}

export const config = {
  AUTH0_DOMAIN: required("AUTH0_DOMAIN"),
  AUTH0_CLIENT_ID: required("AUTH0_CLIENT_ID"),
  AUTH0_CLIENT_SECRET: required("AUTH0_CLIENT_SECRET"),
  AUTH0_AUDIENCE: required("AUTH0_AUDIENCE"),
  AUTH0_SCOPE: process.env.AUTH0_SCOPE || "flags:read",
  API_BASE_URL: process.env.API_BASE_URL || "http://localhost:8080",
  TRAFFIC_INTERVAL_MS: parseInt(process.env.TRAFFIC_INTERVAL_MS || "5000", 10),
  TRAFFIC_CONCURRENCY: parseInt(process.env.TRAFFIC_CONCURRENCY || "3", 10),
};
