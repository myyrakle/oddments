function signup(id, password, attributes) {
    return new Promise((resolve, reject) => {
        userPool.signUp(
            id,
            password,
            Object.entries(attributes).map(
                (e) => new CognitoUserAttribute({ Name: e[0], Value: e[1] })
            ),
            null,
            function (error, result) {
                if (error) {
                    return reject(error);
                } else {
                    return resolve(result);
                }
            }
        );
    });
}

/*
userPool.signUp(
    "sssang97", //username
    "foobar123", //password
    [
        new CognitoUserAttribute({
            Name: "email",
            Value: "sssang97@naver.com",
        }),
        new CognitoUserAttribute({ Name: "nickname", Value: "myyrakle" }),
    ],
    null,
    function (error, result) {
        if (error) {
            console.log(error);
        } else {
            console.log(result);
        }
    }
);
*/
