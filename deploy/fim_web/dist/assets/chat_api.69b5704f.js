import{a8 as a}from"./index.00362c69.js";function s(){return a.get("/api/chat/session",{params:{limit:-1}})}function e(t){return a.get("/api/chat/history",{params:t})}function r(t){return a.delete("/api/chat/chat",{data:{idList:t}})}function o(t){return a.post("/api/chat/user_top",{friendId:t})}export{o as a,e as b,s as c,r as d};
