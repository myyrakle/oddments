/*
npm i --save amazon-cognito-identity-js
*/

const {
    CognitoUserAttribute,
    CognitoUserPool,
    CognitoUser,
    AuthenticationDetails,
} = require("amazon-cognito-identity-js");

const userPool = new CognitoUserPool({
    UserPoolId: "...",
    ClientId: "...",
});

module.exports = userPool
