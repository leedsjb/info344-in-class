FROM node
RUN mkdir /app

# set working directory to newly created app dir
WORKDIR /app

# copy just package.json first so npm install can be run to detect any changes to package.json
COPY package.json /app

# look inside package.json and install all necessary dependencies
RUN npm install

# copy all source code 
COPY . /app

# expose can be used w/ the docker run -P flag
EXPOSE 80

CMD ["node", "index.js"]

# commands
# docker build -t leeds/nodezip .
# docker run -d -p 80:80 --name nodezip leeds/nodezip