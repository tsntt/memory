const Memory = (cfg) => {
  const canvas = document.getElementById(cfg.id)
  const scene = new THREE.Scene();
  const renderer = new THREE.WebGLRenderer(cfg.renderer);
  const camera = new Camera();

  renderer.setPixelRatio(window.devicePixelRatio);
  renderer.setSize(window.innerWidth, window.innerHeight)

  renderer.shadowMap.enabled = true;
  renderer.shadowMap.type = THREE.PCFSoftShadowMap;

  cfg.cards.camera = camera
  cfg.cards.scene = scene

  const cards = Cards(cfg.cards)
  const lights = Lights(scene)

  const geometry = new THREE.PlaneGeometry( 50, 50 );
  const material = new THREE.MeshStandardMaterial( {color: 0xdddddd, side: THREE.DoubleSide} );
  const plane = new THREE.Mesh(geometry, material);
  plane.rotation.x = Math.PI / 2;
  plane.position.x = 5
  plane.position.z = 5
  plane.position.y = -0.02
  plane.receiveShadow = true;
  scene.add( plane );

  function animate() {
    requestAnimationFrame( animate );
    renderer.render( scene, camera.Get() );
  }

  function stopAnimate() {
    cancelAnimationFrame(animate)
  }

  function resize() {
    camera.Resize();
    renderer.setSize( window.innerWidth, window.innerHeight );
  }

  function cardHover(event) {
    cards.raycasterEvent(event).do(card => {
      if (card) {
        const pos = card.position;
        lights[1].position.set(pos.x, 0.1, pos.z);
        lights[1].lookAt(pos.x, 0, pos.z);
        lights[1].updateMatrix();
        lights[1].visible = true;
      } else {
        lights[1].visible = false;
      }
    });
  }

  // put this on window.Memory
  var m = null
  m = {
    New: () => {
      // remove previous things if exists
      stopAnimate();
      window.removeEventListener('resize', resize);
      window.removeEventListener('mousemove', cardHover);
      canvas.innerHTML = '';

      cfg.loading.show();
      
      // new one
      canvas.appendChild(renderer.domElement);
      window.addEventListener('resize', resize);
      window.addEventListener('mousemove', cardHover);
      animate();

      return m
    },
    Cards: {
      select: (fn) => {
        window.addEventListener('mousedown', (event) => { 
          if (event.target == renderer.domElement) {
            cards.raycasterEvent(event).do(card => { 
              fn(cards.id(card)) 
            })
          }
        })
      },
      reveal: (id, content) => {
        cards.Turn(id, content)
      },
      faceDown: (idA, idB) => {
        cards.Turn(idA)
        cards.Turn(idB, '', 0.1)
      },
      paired: (idA, idB, userID) => {
        const elem = document.getElementById(userID).querySelectorAll(cfg.paired.target)[0]
        const bb = elem.getBoundingClientRect()
        const vec2 = new THREE.Vector2(1, 1)
      
        vec2.x = ( bb.x / window.innerWidth ) * 2 - 1;
        vec2.y = - ( bb.y / window.innerHeight ) * 2 + 1;
      
        cards.RemovePaired(idA, idB, vec2, _ => {
          console.log('finished')
        })
      }
    }
  }

  return m
};