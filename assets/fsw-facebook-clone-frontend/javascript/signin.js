let signin_api = "http://localhost/fsw-facebook-clone-backend/php/signin_api.php";

const signIn = async (email, password) => {
    const response = await fetch(signin_api, {
        method: "POST",
        headers: new Headers({
            "Content-Type": "application/json"
        })
        ,body: JSON.stringify({
            email: email,
            password: password
        })
    });
    const json_object = await response.json();
    if (json_object.status == "Logged In"){
        token = json_object.token;
        localStorage.setItem("token", token);
        location.href = "http://localhost/fsw-facebook-clone-frontend/home.html";
        return token;
    }else{
        console.log(json_object.status);
    }
}  

document.getElementById("signin_button").addEventListener("click", function(){
    let email = document.getElementById("email1").value;
    let password = document.getElementById("password1").value;   
    signIn(email, password);
});

