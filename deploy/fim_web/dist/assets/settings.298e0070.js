import{d as y,b as h,r as w,x as C,c as _,g as u,w as t,l as s,a as p,f as a,m as i,i as f,h as Q,K as k,ay as U,E as F,az as x,k as S}from"./index.00362c69.js";const I={class:"group_settings"},T=a("div",{style:{"font-size":"14px",color:"var(--text_color)"}},[a("span",null,"\uFF08\u52FE\u9009\u540E\u5C06\u542F\u7528\u7FA4\u53F7\uFF0C\u4ED6\u4EBA\u53EF\u6309\u7FA4\u53F7\u3001\u7FA4\u540D\u79F0\u641C\u7D22\u672C\u7FA4\uFF09")],-1),z={key:0,class:"verificationQuestion"},N={class:"item"},P=a("span",{class:"label"},"\u95EE\u98981",-1),G={key:0},K=a("span",{class:"label"},"\u7B54\u6848",-1),M={class:"item"},R=a("span",{class:"label"},"\u95EE\u98982",-1),j={key:0},q=a("span",{class:"label"},"\u7B54\u6848",-1),H={class:"item"},J=a("span",{class:"label"},"\u95EE\u98983",-1),L={key:0},O=a("span",{class:"label"},"\u7B54\u6848",-1),W={class:"radio_group"},Z=y({__name:"settings",setup(X){const n=h(),e=w({id:n.groupData.groupId,isSearch:!1,verification:2,isInvite:!1,isTemporarySession:!1,verificationQuestion:{problem1:"",problem2:"",problem3:"",answer1:"",answer2:"",answer3:""}});C(()=>n.groupData,()=>{e.id=n.groupData.groupId,e.isProhibition=n.groupData.isProhibition,e.isSearch=n.groupData.isSearch,e.isInvite=n.groupData.isInvite,e.isTemporarySession=n.groupData.isTemporarySession},{immediate:!0,deep:!0});async function B(){let r=await U(e);if(r.code){F.error(r.msg);return}F.success("\u4FEE\u6539\u7FA4\u8D44\u6599\u6210\u529F")}async function V(){let r=await x(e.id);if(r.code){F.error(r.msg);return}F.success("\u7FA4\u89E3\u6563\u6210\u529F"),S.push({name:"session"})}return(r,o)=>{const v=s("el-checkbox"),m=s("el-form-item"),c=s("el-radio"),b=s("el-radio-group"),d=s("el-input"),E=s("el-button"),A=s("el-popconfirm"),g=s("el-form"),D=s("el-scrollbar");return p(),_("div",I,[u(D,{height:"500px"},{default:t(()=>[u(g,{model:e},{default:t(()=>[u(m,{label:"\u67E5\u627E\u65B9\u5F0F"},{default:t(()=>[a("div",null,[a("div",null,[u(v,{modelValue:e.isSearch,"onUpdate:modelValue":o[0]||(o[0]=l=>e.isSearch=l)},{default:t(()=>[i("\u5141\u8BB8\u88AB\u641C\u7D22")]),_:1},8,["modelValue"])]),T])]),_:1}),u(m,{label:"\u52A0\u7FA4\u65B9\u5F0F"},{default:t(()=>[a("div",null,[a("div",null,[u(b,{class:"radio_group",modelValue:e.verification,"onUpdate:modelValue":o[1]||(o[1]=l=>e.verification=l)},{default:t(()=>[u(c,{value:0},{default:t(()=>[i("\u4E0D\u5141\u8BB8\u4EFB\u4F55\u4EBA\u52A0\u7FA4")]),_:1}),u(c,{value:1},{default:t(()=>[i("\u5141\u8BB8\u4EFB\u4F55\u4EBA\u52A0\u7FA4")]),_:1}),u(c,{value:2},{default:t(()=>[i("\u9700\u8981\u9A8C\u8BC1\u6D88\u606F")]),_:1}),u(c,{value:3},{default:t(()=>[i("\u9700\u8981\u56DE\u7B54\u95EE\u9898")]),_:1}),u(c,{value:4},{default:t(()=>[i("\u9700\u8981\u6B63\u786E\u56DE\u7B54\u95EE\u9898")]),_:1})]),_:1},8,["modelValue"])]),e.verification===3||e.verification===4?(p(),_("div",z,[a("div",N,[a("div",null,[P,u(d,{placeholder:"\u95EE\u98981",modelValue:e.verificationQuestion.problem1,"onUpdate:modelValue":o[2]||(o[2]=l=>e.verificationQuestion.problem1=l)},null,8,["modelValue"])]),e.verification===4?(p(),_("div",G,[K,u(d,{placeholder:"\u7B54\u6848",modelValue:e.verificationQuestion.answer1,"onUpdate:modelValue":o[3]||(o[3]=l=>e.verificationQuestion.answer1=l)},null,8,["modelValue"])])):f("",!0)]),a("div",M,[a("div",null,[R,u(d,{placeholder:"\u95EE\u98982",modelValue:e.verificationQuestion.problem2,"onUpdate:modelValue":o[4]||(o[4]=l=>e.verificationQuestion.problem2=l)},null,8,["modelValue"])]),e.verification===4?(p(),_("div",j,[q,u(d,{placeholder:"\u7B54\u6848",modelValue:e.verificationQuestion.answer2,"onUpdate:modelValue":o[5]||(o[5]=l=>e.verificationQuestion.answer2=l)},null,8,["modelValue"])])):f("",!0)]),a("div",H,[a("div",null,[J,u(d,{placeholder:"\u95EE\u98983",modelValue:e.verificationQuestion.problem3,"onUpdate:modelValue":o[6]||(o[6]=l=>e.verificationQuestion.problem3=l)},null,8,["modelValue"])]),e.verification===4?(p(),_("div",L,[O,u(d,{placeholder:"\u7B54\u6848",modelValue:e.verificationQuestion.answer3,"onUpdate:modelValue":o[7]||(o[7]=l=>e.verificationQuestion.answer3=l)},null,8,["modelValue"])])):f("",!0)])])):f("",!0)])]),_:1}),u(m,{label:"\u9080\u8BF7\u65B9\u5F0F"},{default:t(()=>[u(v,{modelValue:e.isInvite,"onUpdate:modelValue":o[8]||(o[8]=l=>e.isInvite=l)},{default:t(()=>[i("\u5141\u8BB8\u7FA4\u6210\u5458\u9080\u8BF7\u597D\u53CB\u8FDB\u7FA4")]),_:1},8,["modelValue"])]),_:1}),u(m,{label:"\u4F1A\u8BDD\u6743\u9650"},{default:t(()=>[a("div",W,[u(v,{modelValue:e.isTemporarySession,"onUpdate:modelValue":o[9]||(o[9]=l=>e.isTemporarySession=l),disabled:""},{default:t(()=>[i("\u5141\u8BB8\u7FA4\u6210\u5458\u53D1\u8D77\u4E34\u65F6\u4F1A\u8BDD")]),_:1},8,["modelValue"]),u(v,{modelValue:e.isProhibition,"onUpdate:modelValue":o[10]||(o[10]=l=>e.isProhibition=l)},{default:t(()=>[i("\u5168\u5458\u7981\u8A00")]),_:1},8,["modelValue"])])]),_:1}),u(m,{label:"\u5176\u4ED6"},{default:t(()=>[u(E,{size:"small",type:"primary",onClick:B},{default:t(()=>[i("\u4FEE\u6539\u7FA4\u8BBE\u7F6E")]),_:1}),Q(n).groupData.role===1?(p(),k(A,{key:0,title:"\u662F\u5426\u89E3\u6563\u8BE5\u7FA4\uFF1F",onConfirm:V},{reference:t(()=>[u(E,{size:"small",type:"danger"},{default:t(()=>[i("\u89E3\u6563\u8BE5\u7FA4")]),_:1})]),_:1})):f("",!0)]),_:1})]),_:1},8,["model"])]),_:1})])}}});export{Z as default};