import{d as F,u as x,b as h,x as y,c as C,f as a,y as B,h as d,g as n,w as o,l as _,a as l,m as r,K as p,i as m,k as f,N as E,E as v}from"./index.00362c69.js";const b={class:"group_chat_view"},N={class:"group_chat_head"},I={class:"title"},V={class:"icon"},$=a("i",{class:"iconfont icon-more"},null,-1),G={class:"group_chat_main"},S=F({__name:"index",setup(M){const g=x(),t=h();y(()=>g.params,()=>{t.getGroupData(Number(g.params.id))},{immediate:!0,deep:!0});function i(s){f.push({name:s})}async function w(){let s=await E({memberId:t.userInfo.userID,id:t.groupData.groupId});if(s.code){v.error(s.msg);return}v.success("\u9000\u51FA\u7FA4\u804A\u6210\u529F"),f.push({name:"session"})}return(s,e)=>{const u=_("el-dropdown-item"),k=_("el-dropdown-menu"),A=_("el-dropdown"),D=_("router-view");return l(),C("div",b,[a("div",N,[a("div",I,B(d(t).groupData.title),1),a("div",V,[n(A,{trigger:"click"},{dropdown:o(()=>[n(k,null,{default:o(()=>[n(u,{onClick:e[0]||(e[0]=c=>i("session_group_chat"))},{default:o(()=>[r("\u7FA4\u5BF9\u8BDD")]),_:1}),n(u,{onClick:e[1]||(e[1]=c=>i("group_information"))},{default:o(()=>[r("\u7FA4\u8D44\u6599")]),_:1}),d(t).groupData.role!==3?(l(),p(u,{key:0,onClick:e[2]||(e[2]=c=>i("group_member"))},{default:o(()=>[r("\u7FA4\u6210\u5458 ")]),_:1})):m("",!0),d(t).groupData.role!==3?(l(),p(u,{key:1,onClick:e[3]||(e[3]=c=>i("group_settings"))},{default:o(()=>[r("\u7FA4\u8BBE\u7F6E ")]),_:1})):m("",!0),d(t).groupData.role!==1?(l(),p(u,{key:2,style:{color:"red"},onClick:w,divided:""},{default:o(()=>[r(" \u9000\u51FA\u7FA4\u804A ")]),_:1})):m("",!0)]),_:1})]),default:o(()=>[$]),_:1})])]),a("div",G,[n(D)])])}}});export{S as default};
