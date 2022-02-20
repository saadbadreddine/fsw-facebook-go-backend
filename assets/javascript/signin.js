const signIn = async (email, password) => {
    const response = await fetch("http://localhost:8080/signin", {
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
    if (json_object.status == "Logged in"){
        token = json_object.token;
        localStorage.setItem("token", token);
        location.href = "http://localhost:8080/assets/fsw-facebook-clone-frontend/home.html";
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

