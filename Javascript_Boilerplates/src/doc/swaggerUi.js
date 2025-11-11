const swaggerJSDoc = require("swagger-jsdoc");
const swaggerUi = require("swagger-ui-express");

const swaggerDefinition = {
    info: {
        title: "DAWNFOODS",
        version: "1.0.0",
        description: "DAWNFOODS API DOCUMENT",
    },
    host: "672r7zvf6d.execute-api.ap-northeast-2.amazonaws.com/prod/v1",
    basePath: "/",
};

const option = {
    swaggerDefinition,
    apis: [`${__dirname}/route/v1.js`],
};

// swagger-jsdoc 초기화.
const swaggerSpec = swaggerJSDoc(option);

module.exports = {
    serve: swaggerUi.serve,
    setup: swaggerUi.setup(swaggerSpec),
    middleware: swaggerUi.serveFiles(swaggerSpec, option),
    response: (req, res) =>
        res.send(swaggerUi.generateHTML(swaggerSpec, option)),
};

/**
app.use("/doc", user.middleware);
app.get("/doc", user.response);
 */
