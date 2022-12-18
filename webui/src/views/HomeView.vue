<script setup>
import LoadingSpinner from '../components/LoadingSpinner.vue';
import Photo from '../components/Photo.vue'
</script>

<script>
export default {
	components: {
    Photo,
    LoadingSpinner
},
	data: function() {
		return {
			errormsg: null,
			loading: false,
			feed: [],
			nextPageId: 0,
		}
	},
	methods: {
		async getFeed() {
			this.loading = true;
			this.errormsg = null;
			try{
				let response = await this.$axios.get("/feed/");
				this.feed.push(response.data.posts);
				this.nextPageId = response.data.nextFeedPageId;
			} catch(e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
  		},
		async getPhoto(id) {
			this.loading = true;
			this.errormsg = null;
			try{
				let response = await this.$axios.get(`/feed/${id}/`);
				return response.data;
			} catch(e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		}
	},
	mounted() {
		if (this.$profile.identifier==null){
			this.$router.push({path: '/'});
		}

		this.getFeed();
	}
}
</script>

<template>
	<div class="d-flex min-vh-100 w-100 justify-content-center align-items-center" style="background-color: #383838">
		<div class="min-vh-100 w-75" style="background-color: #2e2e2e;">		
			<div class="d-flex flex-column min-vh-75">
				<Photo v-for="photo in feed"/>
			</div>
			<div class="text-white" v-if="nextPageId==0">
				Fine Feed.
			</div>
		</div>
	</div>
</template>

<style>
</style>
