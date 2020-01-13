Setting up a federation of Scrit mints
--------------------------------------

This walkthrough shows how one would set up a federation of 3 Scrit
mints.

First each of the three mints has to generate their mint identity key:

    $ scrit-mint keygen

Afterwards they all display their identity key and send it to the
federation coordinator (one of the mints should manage the
[Codechain](https://github.com/frankbraun/codechain) for the new Scrit
federation).

    $ scrit-mint identity

Typical output looks like this:

    ed25519-fZLPEvdKwhvxU_asrnqbR9t1PV0FukT71f1iwExX_ic

Let's say we have the following three mint identity keys:

    ed25519-vVqGX7eEyH5DNxO_UHm2k8iJAvf-NNv2g1UbZnTnu44
    ed25519-boVnUGMNKkI1Pe72m8Kf_9KljL4DBvsOGxbr1wi9flo
    ed25519-er0Phn1PjBzbz3gBUEbFQUIbexZxufELZyzCyfT4A5U

Now we can setup the federation (2-of-3):

        $ scrit-gov start -m 2 -n 3 d25519-vVqGX7eEyH5DNxO_UHm2k8iJAvf-NNv2g1UbZnTnu44 ed25519-boVnUGMNKkI1Pe72m8Kf_9KljL4DBvsOGxbr1wi9flo ed25519-er0Phn1PjBzbz3gBUEbFQUIbexZxufELZyzCyfT4A5U

To be continued...
