<script>
export default {
	props: ['post'],
	emits: ['photoDeleted'],
    data: function() {
		return {
			postDetails: Object.assign({}, this.post),
			errormsg: null,
			loading: null,
			blob: null,
			blobUrl: null,
			showCommentStatus: false,
			liked: false,
			comments: [],
			nextCommentPageId: 0,
			newComment: null,
		}
	},
	methods: {
		async getPhoto() {
			this.loading = true;
			this.errormsg = null;
			try{
				let response = await this.$axios.get(`/feed/${this.post.photo}/`, { responseType: 'blob'});
				this.blob = response.data
				this.blobUrl = URL.createObjectURL(this.blob);
				document.getElementById(`${this.post.photo}`).src = this.blobUrl;
			} catch(e) {
				this.errormsg = e.response.data.message;
			}
			this.loading = false;
  		},
		async getStatus() {
			this.loading = true;
			this.errormsg = null;
			// status
			try{
				let response = await this.$axios.get(`/feed/${this.post.photo}/likes/${this.$profile.username}`);
				this.liked = response.data.status;
			} catch(e) {
				this.errormsg = e.response.data.message;
			}
			this.loading = false;
  		},
		async setLike() {
			this.loading = true;
			this.errormsg = null;
			try{
				let response = await this.$axios.put(`/feed/${this.post.photo}/likes/${this.$profile.username}`);
				this.liked = true;
				this.postDetails.numberOfLikes += 1;
			} catch(e) {
				if (e.response.status == 409) {
					this.errormsg = e.response.data.message;
					this.liked = true;
				}
				
			}
			this.loading = false;
  		},
		async removeLike() {
			this.loading = true;
			this.errormsg = null;
			try{
				let response = await this.$axios.delete(`/feed/${this.post.photo}/likes/${this.$profile.username}`);
				this.liked = false;
				this.postDetails.numberOfLikes -= 1;
			} catch(e) {
				if (e.response.status == 409) {
					this.errormsg = e.response.data.message;
					this.liked = false;
				}
			}
			this.loading = false;
  		},
		showComment() {
			this.showCommentStatus = !this.showCommentStatus;
		},
		async deleteComment(commentId) {
            this.loading = true;
			this.errormsg = null;
			try{
				let response = await this.$axios.delete(`/feed/${this.post.photo}/comments/${this.$profile.username}?idCommento=${commentId}`);
				this.postDetails.numberOfComments -= 1;
				for (var i=0; i<this.comments.length; i++){
					if (this.comments[i].id === commentId){
						this.comments.splice(i,1)
						break;
					}
				}
			} catch(e) {
				this.errormsg = e.response.data.message;
			}
			this.getComments();
        },
		async sendComment(commentId) {
			if (this.newComment != null) {
				this.loading = true;
				this.errormsg = null;
				try{
					let response = await this.$axios.post(`/feed/${this.post.photo}/comments/${this.$profile.username}`, {text: this.newComment});
					this.postDetails.numberOfComments += 1;
					this.newComment = null
				} catch(e) {
					this.errormsg = e.response.data.message;
				}
				this.getComments();
			}
        },
		async getComments() {
			this.loading = true;
			this.errormsg = null;
			try{
				const nextPage = this.nextCommentPageId;
				let response = nextPage == 0?  	await this.$axios.get(`/feed/${this.post.photo}/comments/`) : 
												await this.$axios.get(`/feed/${this.post.photo}/comments/?pageId=${this.nextCommentPageId}`) ;
				if (!(response.status === 204)) {
					nextPage == 0 ? this.comments = response.data.comments : this.comments.push(...response.data.comments);
					this.nextCommentPageId = response.data.nextCommentPageId;
				}
			} catch(e) {
				this.errormsg = e.response.data.message;
			}
			this.loading = false;
		},
		visitProfile(username) {
			this.$router.push({ path: `/profiles/${username}` });
		},
		async deletePhoto() {
			this.loading = true;
			this.errormsg = null;
			try{
				let response = await this.$axios.delete(`/profiles/${this.$profile.username}/photos/${this.post.photo}`);
				this.$emit('photoDeleted', this.post.photo);
			} catch(e) {
				this.errormsg = e.response.data.message;
			}
			this.loading = false;
		}
	},
	mounted() {
		this.getPhoto();
		this.getStatus();
		this.getComments();
	},
}
</script>

<template>
	<LoadingSpinner :loading="this.loading" class="pb-2"/>
	<ErrorMsg :msg="this.errormsg"/>
	<div class="card" style="background-color: #212121;">
        <div class="PhotoHeader">
			<div class="username fw-bold fs-3" @click="visitProfile(this.post.username)">
				{{ this.postDetails.username }}
			</div>
			<div class="d-flex flex-row justify-content-end align-items-center w-75">
				<div class="data fs-6"> 
					{{ this.postDetails.data }}
				</div>
				<div v-if="this.postDetails.username==this.$profile.username" class="elimina fs-6">
					<i class="icon bi bi-trash3-fill text-danger fs-3" @click="deletePhoto"></i>
				</div>
			</div>
		</div>		
		<div class="d-flex justify-content-center align-items-center">
			<img :id="this.postDetails.photo" class="image">
		</div>
		<div class="Buttons">
			<div class="d-flex justify-content-start align-items-center min-vw-25 max-vw-25">
				<div class="text-white fs-bold fs-4 pe-2">
					{{ this.postDetails.numberOfLikes }}
				</div>
				<i v-if="this.liked" class="icon bi bi-heart-fill" @click="removeLike" style="color:red"></i>
				<i v-else class="icon bi bi-heartbreak" @click="setLike" style="color:red"></i>
			</div>
			<div class="d-flex justify-content-center align-items-center min-vw-25 max-vw-25">
				<div class="text-white fs-bold fs-4 pe-2">
					{{ this.postDetails.numberOfComments }}
				</div>
				<i class="icon bi bi-chat-left-text" @click="showComment"></i>
			</div>
		</div>	
		<div class="commentAreaInput">
			<textarea v-model="this.newComment" class="textArea border border-2 border-success" placeholder="Commento" style="background-color: #212121;"></textarea>
			<i class="icon bi bi-send-fill text-success fs-4" @click="sendComment"></i>
		</div>
		<div v-if="this.showCommentStatus" class="commentArea">
			<div v-for="comment in this.comments" v-bind:key="comment" class="comment pb-2">
				<div class="dettagli">
					<div class="fw-bolder fs-5 text-white pe-4" @click="visitProfile(comment.username)">
						{{ comment.username }}
					</div>
					<div class="fw-normal fs-6 text-white">
						{{ comment.text }}
					</div>
				</div>
				<div class="dettagli">
					<div class="text-bold text-success fs-6 pe-4">
						{{ comment.data }}
					</div>
					<div v-if="((comment.idUser===this.$profile.identifier) || (this.postDetails.username===this.$profile.username))" class="text-decoration-underline text-success fs-6" @click="deleteComment(comment.id)">
						Elimina
					</div>
				</div>
			</div>
			<div v-if="this.nextCommentPageId!=0" class="p-4">
				<button type="button" class="btn btn-outline-success text-white fw-bolder rounded-pill fs-5" style="width: 150px" @click="getComments">...</button>
			</div>
		</div>
    </div>
	<div class="p-2"/>
</template>

<style scoped>
	.card {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		padding-top: 2rem;
		width: 60%;
		border-radius: 1rem;
	}
	.PhotoHeader {
		display: flex;
		flex-direction: row;
		justify-content: space-between;
		align-items: center;
		width: 90%;
	}
	.username {
		display: flex;
		justify-content: start;
		align-items: center;
		width: 90%;
		padding-bottom: 1rem;
		color: white;
	}
	.image {
		object-fit: contain;
		width: 90%;
		padding-bottom: 0.5rem;
	}
	.data {
		display: flex;
		justify-content: right;
		align-items: center;
		width: 90%;
		padding-bottom: 1rem;
		color: white;
	}

	.elimina {
		display: flex;
		justify-content: right;
		align-items: center;
		padding-bottom: 1rem;
		padding-left: 1rem;
		color: white;
	}
	.Buttons {
		display: flex;
		justify-content: space-between;
		align-items: center;
		width: 90%;
	}

	.commentAreaInput {
		display: flex;
		justify-content: space-between;
		align-items: center;
		width: 90%;
		height: 10%;
		padding-bottom: 1rem;
	}

	.textArea {
		min-width: 92%;
		border-radius: 2rem;
		height: 2.25rem;
		color: white;
		text-indent: 10px;
	}

	.commentArea {
		display: flex;
		flex-direction: column;;
		justify-content: space-between;
		align-items: center;
		width: 90%;
		overflow: scroll;
		max-height: 150px;
		padding-bottom: 1rem;
	}

	.comment {
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        width: 100%;
    }

    .dettagli {
        display: flex;
        flex-direction: row;
        justify-content: start;
        align-items: center;
        width: 100%;
    }

</style>

