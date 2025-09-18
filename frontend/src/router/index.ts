import { createRouter, createWebHistory, RouteRecordRaw  } from "vue-router";
import Home from "../pages/Home.vue";
import DocumentsPage from "../pages/DocumentsPage.vue";
import CatText from "../pages/CatText.vue";
import JigsawPage from "../pages/JigsawPage.vue";
import FeedbackPage from "../pages/FeedbackPage.vue";
import CatGame from "../pages/CatGame.vue";
// import DocumentsPageAlls from "../pages/DocumentsPageAlls.vue";
import DogAll from '../pages/DogAll.vue'

const routes = [
  { path: "/", component: Home },
  { path: "/DogPuzzle", name: "DocumentsPage", component: DocumentsPage },
  { path: "/CatText", name: "CatText", component: CatText },
  { path: "/CatGame", name: "CatGame", component: CatGame },
  { path: "/PolaJigsaw", name: "JigsawPage", component: JigsawPage },
  { path: "/Feedback", name: "FeedbackPage", component: FeedbackPage },
  { path: "/party", name: "DocumentsPageAlls", component: () => import("../pages/DocumentsPageAlls.vue") },
{ path: "/party/:code", name: "DocumentsPageRoom", component: () => import("../pages/DocumentsPageAlls.vue"), props: true },
 { path: '/dog', name: 'DogAll', component: DogAll },

];

export default createRouter({
  history: createWebHistory(),
  routes,
});
