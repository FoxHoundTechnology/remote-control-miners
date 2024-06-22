# Use an official Node.js runtime as the base image
FROM node:20-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the package.json and package-lock.json files to the container
COPY package.json package-lock.json* ./

# Install dependencies
RUN npm install

# Copy the rest of the application files to the container
COPY . .

# Build the Next.js project
RUN npm run build

# Expose the port that the application will run on
EXPOSE 3000

# Define the command to run the application in production mode
CMD ["npm", "start"]