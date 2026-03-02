.PHONY: addlicense

ADDLICENSE ?= addlicense
COPYRIGHT_HOLDER ?= Woodstock K.K.
LICENSE_TYPE ?= mit

addlicense:
	$(ADDLICENSE) -c "$(COPYRIGHT_HOLDER)" -l "$(LICENSE_TYPE)" .
