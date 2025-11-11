function login(id, password) {
    const user = new CognitoUser({ Username: id, Pool: userPool });

    return new Promise((resolve, reject) => {
        user.authenticateUser(
            new AuthenticationDetails({
                Username: id,
                Password: password,
            }),
            {
                onSuccess: (result) => {
                    resolve({ userSession: result, user: user });
                },
                onFailure: (error) => {
                    console.log(error);
                },
            }
        );
    });
}

/*
const user = new CognitoUser({ Username: "sssang97", Pool: userPool });
user.authenticateUser(
    new AuthenticationDetails({ Username: "sssang97", Password: "password" }),
    {
        onSuccess: (result) => {
            console.log("로그인됨");

            console.log(result);

            const token = result.getAccessToken().getJwtToken();
            console.log("token:", token);

            user.getUserAttributes((error, result) => {
                if (error) {
                    console.log(error);
                } else {
                    console.log("attr:", result);
                }
            });
        },
        onFailure: (error) => {
            console.log(error);
        },
    }
);
*/
