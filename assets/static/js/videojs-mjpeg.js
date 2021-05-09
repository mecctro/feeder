(function (root, factory) {
  if(typeof exports==='object' && typeof module!=='undefined') {
    module.exports = factory(require('video.js'));
  } else if(typeof define === 'function' && define.amd) {
    define(['videojs'], function(videojs){
      return (root.Mjpeg = factory(videojs));
    });
  } else {
    root.Mjpeg = factory(root.videojs);
  }
}(this, function(videojs) {
  'use strict';

  var Tech = videojs.getTech('Tech');
  var Loader;

  var Mjpeg = videojs.extend(Tech, {
    constructor: function(options, ready) {
      Tech.call(this, options, ready);

      this.setSrc(options.source, true);

      // Set the vjs-youtube class to the player
      // Parent is not set yet so we have to wait a tick
      setTimeout(function() {
        this.el_.parentNode.className += ' vjs-mjpeg';
        this.loadFrame();
        this.triggerReady();
      }.bind(this));
    },

    dispose: function() {
      this.trigger('stop');
      cleartInterval(Loader);

      this.el_.parentNode.className = this.el_.parentNode.className
        .replace(' vjs-mjpeg', '');
      this.el_.parentNode.removeChild(this.el_);

      Tech.prototype.dispose.call(this);
    },

    currentTime: function() {
      return 0;
    },

    duration: function() {
      return 0;
    },

    seeking: function() {
      return false;
    },

    seekable: function () {
      return false;
    },

    volume: function() {
      return false;
    },

    createEl: function() {
      //console.log(this.options_)
      var canvas = document.createElement('canvas');
      canvas.setAttribute('id', this.options_.techId);
      canvas.setAttribute('style', 'max-width:100%;max-height:100%;height:100%;');
      //div.setAttribute('src', this.options_.source.src);

      var div = document.createElement('div');
      div.setAttribute('class', 'vjs-tech');
      div.setAttribute('style', 'text-align:center;');
      div.appendChild(canvas);

      return div;
    },

    src: function(src) {
      if (src) {
        this.setSrc({ src: src });
      }

      return this.source;
    },

    setSrc: function(source) {
      if (!source || !source.src) {
        return;
      }

      delete this.errorNumber;
      this.source = source;
    },

    autoplay: function() {
      //return this.options_.autoplay;
    },

    setAutoplay: function(val) {
      //this.options_.autoplay = val;
    },

    play: function() {
      Loader = window.setInterval(this.loadFrame.bind(this), this.options_.source.rate);
      this.trigger('play');

      this.el_.parentNode.querySelector('.vjs-progress-control').setAttribute('style', 'display:none;');
      this.el_.parentNode.querySelector('.vjs-duration').setAttribute('style', 'display:none;');
      this.el_.parentNode.querySelector('.vjs-current-time').setAttribute('style', 'display:none;');
      this.el_.parentNode.querySelector('.vjs-remaining-time').setAttribute('style', 'display:none;');
      this.el_.parentNode.querySelector('.vjs-volume-menu-button').setAttribute('style', 'display:none;');
      this.el_.parentNode.querySelector('.vjs-live-control').setAttribute('style', 'display:inherit !important;');


      return true;
    },

    pause: function() {
      //this.el_.parentNode.parentNode.className.replace(' vjs-playing', ' vjs-paused');
      //this.el_.parentNode.querySelector('.vjs-live-control').className.replace(' vjs-playing', ' vjs-paused');
      this.trigger('pause');
      cleartInterval(Loader);

      return true;
    },

    paused: function() {

      return true;
    },

    loadFrame: function() {
      var canvas = document.getElementById(this.options_.techId);
      var ctx = canvas.getContext('2d');
      var img = new Image();
      img.src = this.options_.source.src;
      img.onload = function() {
        canvas.width = this.width;
        canvas.height = this.height;
        //alert(this.width + ' ' + this.height);
        ctx.drawImage(img, 0, 0);
      };
      console.log('loading frame...')
    }

  });

  Mjpeg.isSupported = function() {
    return true;
  };

  Mjpeg.canPlaySource = function(srcObj){
    if (srcObj.type === 'multipart/x-mixed-replace') {
      // TODO: allow codec info and check browser support
      return 'maybe';
    } else {
      return '';
    }
  };

  videojs.registerTech('Mjpeg', Mjpeg);

}));
