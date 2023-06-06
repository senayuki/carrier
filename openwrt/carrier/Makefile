include $(TOPDIR)/rules.mk

PKG_NAME:=carrier
PKG_VERSION:=0.0.1
PKG_RELEASE:=$(AUTORELEASE)
PKG_MAINTAINER:=Yuki Sena <admin@senayuki.moe>

PKG_SOURCE:=$(PKG_NAME)-$(PKG_VERSION).tar.gz
PKG_SOURCE_URL:=https://codeload.github.com/senayuki/carrier/tar.gz/v${PKG_VERSION}?
PKG_HASH:=2d713cc1c1a9fef3a6b06225004ff6ab35b859f7ca6d58b5a06db1711431098e

PKG_BUILD_DIR:=$(BUILD_DIR)/$(PKG_NAME)-$(PKG_VERSION)
PKG_BUILD_DEPENDS:=golang/host
PKG_BUILD_PARALLEL:=1
PKG_USE_MIPS16:=0

GO_PKG:=github.com/senayuki/carrier
GO_PKG_BUILD_PKG:=$$(GO_PKG)
GO_PKG_LDFLAGS:=-s -w
COMPILE_TIME:= $(shell TZ=UTC-8 date '+%Y-%m-%d %H:%M:%S')
GO_PKG_LDFLAGS+= \
	-X '$(GO_PKG)/version.Version=$(PKG_VERSION)' \
	-X '$(GO_PKG)/version.BuildTime=$(COMPILE_TIME)' \
	-X '$(GO_PKG)/version.ExGoVersionInfo=$(GO_ARM) $(GO_MIPS)$(GO_MIPS64)'

include $(INCLUDE_DIR)/package.mk
include $(TOPDIR)/feeds/packages/lang/golang/golang-package.mk

define Package/carrier
	SECTION:=net
	CATEGORY:=Network
	SUBMENU:=Web Servers/Proxies
	TITLE:=A simple reverse proxy
	URL:=https://github.com/senayuki/carrier.git
	DEPENDS:=$(GO_ARCH_DEPENDS)
endef

define Build/Compile
	$(call GoPackage/Build/Configure)
	$(call GoPackage/Build/Compile)
endef

define Package/carrier/install
	$(call GoPackage/Package/Install/Bin,$(PKG_INSTALL_DIR))

	$(INSTALL_DIR) $(1)/usr/bin
	$(INSTALL_BIN) $(PKG_INSTALL_DIR)/usr/bin/carrier $(1)/usr/bin/carrier

	$(INSTALL_DIR) $(1)/usr/share/carrier
	
endef

$(eval $(call GoBinPackage,carrier))
$(eval $(call BuildPackage,carrier))