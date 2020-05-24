$(document).ready(function(){
	$('.imgal-img').click(function(){
		let imageSrc = $(this).attr("src");
		let imageAlt = $(this).attr("alt");

		$('.imgal-container').hide();

		$('body').append(
			'<div class="imgal-modal">'+
			'<span id="imgal-modal-close"">X</span>'+
			'<img src="' + imageSrc + '" alt="' + imageAlt + '" class="imgal-modal-img"></img>'+
			'</div'
		).hide().show('fast');

		$('#imgal-modal-close').click(function(){
			$('.imgal-container').show();
			$('.imgal-modal').hide('fast', function(){
				$(this).remove();
			});
		});
	});
});
