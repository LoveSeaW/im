import{d as C,e as n,c as r,K as E,y as g,g as d,w as I,I as w,E as _,l as p,a,h as x,L as F}from"./index.00362c69.js";const h={class:"fim_friend_notice"},D={key:1},S=C({__name:"fim_friend_notice",props:{friendId:{},modelValue:{}},emits:["update:modelValue"],setup(f,{emit:m}){const l=f,v=m,u=n(!1),s=n(),i=n("");function V(o){u.value=!0,t.value=o,i.value=o,setTimeout(()=>{var e;(e=s.value)==null||e.focus()},100)}const t=n("");async function k(){if(u.value=!1,t.value!==i.value){let o=await w(l.friendId,t.value);if(o.code){_.error(o.msg);return}_.success("\u66F4\u65B0\u597D\u53CB\u5907\u6CE8\u6210\u529F"),v("update:modelValue",t.value)}}return(o,e)=>{const y=p("el-input"),B=p("el-icon");return a(),r("div",h,[u.value?(a(),E(y,{key:0,modelValue:t.value,"onUpdate:modelValue":e[0]||(e[0]=c=>t.value=c),onBlur:k,ref_key:"elIpt",ref:s,placeholder:"\u4FEE\u6539\u597D\u53CB\u5907\u6CE8"},null,8,["modelValue"])):(a(),r("span",D,g(l.modelValue),1)),d(B,{onClick:e[1]||(e[1]=c=>V(l.modelValue)),size:18},{default:I(()=>[d(x(F))]),_:1})])}}});export{S as _};
