const addFriend = async (sender, receiver) => {
  const response = await fetch("http://localhost/fsw-facebook-clone-backend/php/addfriend_api.php", {
    method: "POST",
    headers: new Headers({
      "Content-Type": "application/json",
    }),
    body: JSON.stringify({
      sender: sender,
      receiver: receiver,
    }),
  });
  const data = await response.json();
  return data;
};

const removeFriend = async (sender, receiver) => {
  const response = await fetch("http://localhost/fsw-facebook-clone-backend/php/removefriend_api.php", {
    method: "POST",
    headers: new Headers({
      "Content-Type": "application/json",
    }),
    body: JSON.stringify({
      sender: sender,
      receiver: receiver,
    }),
  });
  const data = await response.json();
  return data;
};

const blockFriend = async (sender, receiver) => {
  const response = await fetch("http://localhost/fsw-facebook-clone-backend/php/blockfriend_api.php", {
    method: "POST",
    headers: new Headers({
      "Content-Type": "application/json",
    }),
    body: JSON.stringify({
      sender: sender,
      receiver: receiver,
    }),
  });
  const data = await response.json();
  return data;
};

const unblockFriend = async (sender, receiver) => {
  const response = await fetch("http://localhost/fsw-facebook-clone-backend/php/unblockfriend_api.php", {
    method: "POST",
    headers: new Headers({
      "Content-Type": "application/json",
    }),
    body: JSON.stringify({
      sender: sender,
      receiver: receiver,
    }),
  });
  const data = await response.json();
  return data;
};

const acceptRequest = async (sender, receiver) => {
  const response = await fetch("http://localhost:8080/acceptfriendrequest", {
    method: "POST",
    headers: new Headers({
      "Content-Type": "application/json",
    }),
    body: JSON.stringify({
      sender: sender,
      receiver: receiver,
    }),
  });
  const data = await response.json();
  return data;
};

const rejectRequest = async (sender, receiver) => {
  const response = await fetch("http://localhost/fsw-facebook-clone-backend/php/rejectrequest_api.php", {
    method: "POST",
    headers: new Headers({
      "Content-Type": "application/json",
    }),
    body: JSON.stringify({
      sender: sender,
      receiver: receiver,
    }),
  });
  const data = await response.json();
  return data;
};

const getfriendRequests = async (token) => {
  const response = await fetch("http://localhost:8080/getfriendrequests", {
    method: "POST",
    headers: new Headers({
      "Content-Type": "application/json",
    }),
    body: JSON.stringify({
      token: token,
    }),
  });
  const data = await response.json();
  document.getElementById("column-two").innerHTML = null;
  data.forEach((element) => {
    first_name = element.first_name.charAt(0).toUpperCase() + element.first_name.slice(1).toLowerCase();
    last_name = element.last_name.charAt(0).toUpperCase() + element.last_name.slice(1).toLowerCase();
    fullname = `${first_name} ${last_name}`;
    document.getElementById("column-two").innerHTML += `
            <div class="friend-container" id="${element.id}">
                <div class="friend-row">
                    <img class="friend-img" src="">
                    <span class="fullname" id="full-name">${fullname}</span>
                </div>
                <div class="button-row">
                    <button class="accept-btn" id="${element.id}-accept">Accept</button>
                    <button class="reject-btn" id="${element.id}-reject">Reject</button>
                </div>
            </div>`;
  });
  return data;
};

const getFriends = async (token) => {
  const response = await fetch("http://localhost:8080/getfriends", {
    method: "POST",
    headers: new Headers({
      "Content-Type": "application/json",
    }),
    body: JSON.stringify({
      token: token,
    }),
  });
  const data = await response.json();
  document.getElementById("column-two").innerHTML = null;
  console.log(data);
  data.forEach((element) => {
    first_name = element.first_name.charAt(0).toUpperCase() + element.first_name.slice(1).toLowerCase();
    last_name = element.last_name.charAt(0).toUpperCase() + element.last_name.slice(1).toLowerCase();
    fullname = `${first_name} ${last_name}`;
    document.getElementById("column-two").innerHTML += `
            <div class="friend-container" id="${element.id}">
                <div class="friend-row">
                    <img class="friend-img" src="">
                    <span class="fullname" id="full-name">${fullname}</span>
                </div>
                <div class="button-row">
                    <button class="block-btn" id="${element.id}-block">Block</button>
                    <button class="remove-btn" id="${element.id}-remove">Remove</button>
                </div>
            </div>`;
  });
  return data;
};

const getblockedUsers = async (token) => {
  const response = await fetch("http://localhost:8080/getblockedusers", {
    method: "POST",
    headers: new Headers({
      "Content-Type": "application/json",
    }),
    body: JSON.stringify({
      token: token,
    }),
  });

  const data = await response.json();
  document.getElementById("column-two").innerHTML = null;
  data.forEach((element) => {
    first_name = element.first_name.charAt(0).toUpperCase() + element.first_name.slice(1).toLowerCase();
    last_name = element.last_name.charAt(0).toUpperCase() + element.last_name.slice(1).toLowerCase();
    fullname = `${first_name} ${last_name}`;
    document.getElementById("column-two").innerHTML += `
            <div class="friend-container" id="${element.id}">
                <div class="friend-row">
                    <img class="friend-img" src="">
                    <span class="fullname" id="full-name">${fullname}</span>
                </div>
                <div class="button-row">
                    <button class="unblock-btn" id="${element.id}-unblock">Unblock</button>
                </div>
            </div>`;
  });
  return data;
};

const getPosts = async (token) => {
  const response = await fetch("http://localhost:8080/getposts", {
    method: "POST",
    headers: new Headers({
      "Content-Type": "application/json",
    }),
    body: JSON.stringify({
      token: token,
    }),
  });

  const data = await response.json();
  document.getElementById("column-two").innerHTML = null;
  data.forEach((element) => {
    post = element.post;
    first_name = element.first_name.charAt(0).toUpperCase() + element.first_name.slice(1).toLowerCase();
    last_name = element.last_name.charAt(0).toUpperCase() + element.last_name.slice(1).toLowerCase();
    fullname = `${first_name} ${last_name}`;

    document.getElementById("column-two").innerHTML += `
            <div class="post-container">
              <div class="user-row">
                    <img class="friend-img" src="">
                    <span class="fullname" id="full-name">${fullname}</span>
                </div>
                <div class="post-row">
                    <span class="post-text" id="post">${post}</span>
                
                 <div class="time-row"></div>`;
  });
  return data;
};

const getData = async (token) => {
  const response = await fetch("http://localhost:8080/getdata", {
    method: "POST",
    headers: new Headers({
      "Content-Type": "application/json",
    }),
    body: JSON.stringify({
      token: token,
    }),
  });
  const json_object = await response.json();
  //if (json_object.status == "User not found") {
  //  console.log(json_object.status);
  //} else {
  return json_object;
  //}
};
const searchUsers = async (sender, name) => {
  const response = await fetch("http://localhost/fsw-facebook-clone-backend/php/searchforusers_api.php", {
    method: "POST",
    headers: new Headers({
      "Content-Type": "application/json",
    }),
    body: JSON.stringify({
      sender: sender,
      name: name,
    }),
  });
  const data = await response.json();
  document.getElementById("column-two").innerHTML = null;
  console.log(data);
  data.forEach((element) => {
    first_name = element.first_name.charAt(0).toUpperCase() + element.first_name.slice(1).toLowerCase();
    last_name = element.last_name.charAt(0).toUpperCase() + element.last_name.slice(1).toLowerCase();
    fullname = `${first_name} ${last_name}`;
    document.getElementById("column-two").innerHTML += `
            <div class="friend-container" id="${element.id}">
                <div class="friend-row">
                    <img class="friend-img" src=${element.picture}>
                    <span class="fullname" id="full-name">${fullname}</span>
                </div>
                <div class="button-row">
                    <button class="add-btn" id="${element.id}-add">Add</button>
                </div>
            </div>`;
  });
  return data;
};

const updatePic = async (sender, picture) => {
  const response = await fetch("http://localhost/fsw-facebook-clone-backend/php/profilepic_api.php", {
    method: "POST",
    headers: new Headers({
      "Content-Type": "application/json",
    }),
    body: JSON.stringify({
      sender: sender,
      picture: picture,
    }),
  });
  const data = await response.json();
  return data;
};
