package main

import (
	"image/color"

	"nftsiren/cmd/nftsiren/alerts"
	"nftsiren/cmd/nftsiren/assets"
	"nftsiren/cmd/nftsiren/widgets"
	"nftsiren/pkg/bench"
	"nftsiren/pkg/nft"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

// Theme holds all styling preferences our app is using
// Theme is not safe for concurrent work, only use it in the frame thread
type Theme struct {
	*material.Theme

	// Colors
	DarkerBg  color.NRGBA
	LighterBg color.NRGBA
	DarkerFg  color.NRGBA
	LighterFg color.NRGBA

	LowImpFg    color.NRGBA
	MediumImpFg color.NRGBA
	HighImpFg   color.NRGBA // Still lower than actual fg

	Link    color.NRGBA
	Success color.NRGBA
	Error   color.NRGBA

	// Units
	CornerRadius    unit.Dp // For all rectangles
	BorderWidth     unit.Dp
	ShadowElevation unit.Dp

	IconSize unit.Dp

	SmallSpace  unit.Dp
	MediumSpace unit.Dp
	LargeSpace  unit.Dp

	SmallInset  layout.Inset
	MediumInset layout.Inset
	LargeInset  layout.Inset

	SmallVSpacer  layout.Spacer
	MediumVSpacer layout.Spacer
	LargeVSpacer  layout.Spacer

	SmallHSpacer  layout.Spacer
	MediumHSpacer layout.Spacer
	LargeHSpacer  layout.Spacer

	// Custom icons
	SirenIcon     *widgets.Icon // Our main icon
	EthereumIcon  *widgets.Icon
	SolanaIcon    *widgets.Icon
	OpenseaIcon   *widgets.Icon
	LooksrareIcon *widgets.Icon
	MagicedenIcon *widgets.Icon

	// Material icons
	HomeIcon        *widgets.Icon
	CollectionsIcon *widgets.Icon
	MintsIcon       *widgets.Icon
	SettingsIcon    *widgets.Icon
	GasIcon         *widgets.Icon
	AlarmIcon       *widgets.Icon
	NewAlarmIcon    *widgets.Icon
	DeleteIcon      *widgets.Icon
	GoBackIcon      *widgets.Icon
	LoopIcon        *widgets.Icon
	BrokenIcon      *widgets.Icon
	FilterIcon      *widgets.Icon
	LogoutIcon      *widgets.Icon
	QuitIcon        *widgets.Icon
}

func DefaultTheme() *Theme {
	defer bench.Begin()()

	mt := material.NewTheme()

	// Main colors
	mt.Bg = rgb(0x212121)
	mt.Fg = rgb(0xeeeeee)
	mt.ContrastBg = rgb(0xFF9800)
	mt.ContrastFg = rgb(0x1A1A1A)

	// TextSize is 16 by default if we don't change it
	var s = unit.Dp(mt.TextSize) / 4 // 4
	var m = s * 2                    // 8
	var l = m * 2                    // 16

	theme := &Theme{
		Theme: mt,

		// Colors
		DarkerBg:  widgets.Darker(mt.Bg),
		LighterBg: widgets.Lighter(mt.Bg),
		DarkerFg:  widgets.Darker(mt.Fg),
		LighterFg: widgets.Lighter(mt.Fg),

		LowImpFg:    widgets.MulAlpha(mt.Fg, 0.25),
		MediumImpFg: widgets.MulAlpha(mt.Fg, 0.50),
		HighImpFg:   widgets.MulAlpha(mt.Fg, 0.75),

		Link:    rgb(0xa0ff10),
		Success: rgb(0xa0ff10),
		Error:   rgb(0xff1010),

		// Units
		CornerRadius:    4,
		BorderWidth:     1,
		ShadowElevation: 2,

		IconSize: unit.Dp(mt.TextSize * 1.5), // 24

		SmallSpace:  s,
		MediumSpace: m,
		LargeSpace:  l,

		SmallInset:  layout.UniformInset(s),
		MediumInset: layout.UniformInset(m),
		LargeInset:  layout.UniformInset(l),

		SmallVSpacer:  layout.Spacer{Height: s},
		MediumVSpacer: layout.Spacer{Height: m},
		LargeVSpacer:  layout.Spacer{Height: l},

		SmallHSpacer:  layout.Spacer{Width: s},
		MediumHSpacer: layout.Spacer{Width: m},
		LargeHSpacer:  layout.Spacer{Width: l},

		// Custom icons
		SirenIcon:     widgets.NewIconFromImage(assets.NftsirenLogo),
		EthereumIcon:  widgets.NewIconFromImage(assets.EthereumLogo),
		SolanaIcon:    widgets.NewIconFromImage(assets.SolanaLogo),
		OpenseaIcon:   widgets.NewIconFromImage(assets.OpenseaLogo),
		LooksrareIcon: widgets.NewIconFromImage(assets.LooksrareLogo),
		MagicedenIcon: widgets.NewIconFromImage(assets.MagicedenLogo),

		// Iconvg icons
		HomeIcon:        widgets.NewIconFromIconVG(icons.ActionHome),
		CollectionsIcon: widgets.NewIconFromIconVG(icons.DeviceWidgets),
		MintsIcon:       widgets.NewIconFromIconVG(icons.ActionEvent),
		SettingsIcon:    widgets.NewIconFromIconVG(icons.ActionSettings),
		GasIcon:         widgets.NewIconFromIconVG(icons.MapsLocalGasStation),
		AlarmIcon:       widgets.NewIconFromIconVG(icons.ActionAlarm),
		NewAlarmIcon:    widgets.NewIconFromIconVG(icons.ActionAlarmAdd),
		DeleteIcon:      widgets.NewIconFromIconVG(icons.ActionDelete),
		GoBackIcon:      widgets.NewIconFromIconVG(icons.NavigationArrowBack),
		LoopIcon:        widgets.NewIconFromIconVG(icons.AVRepeat),
		BrokenIcon:      widgets.NewIconFromIconVG(icons.ImageBrokenImage),
		FilterIcon:      widgets.NewIconFromIconVG(icons.ContentFilterList),
		LogoutIcon:      widgets.NewIconFromIconVG(icons.ActionExitToApp),
		QuitIcon:        widgets.NewIconFromIconVG(icons.ActionPowerSettingsNew),
	}

	return theme
}

// Use this instead of accessing directly
func (theme *Theme) Material() *material.Theme {
	return theme.Theme
}

func (theme *Theme) MarketplaceLogo(market nft.Marketplace) *widgets.Icon {
	switch market {
	case nft.Opensea:
		return theme.OpenseaIcon
	case nft.Looksrare:
		return theme.LooksrareIcon
	case nft.Magiceden:
		return theme.MagicedenIcon
	}
	// This shouldn't happen
	return theme.BrokenIcon
}

func (theme *Theme) AlertIcon(alert alerts.Alert) *widgets.Icon {
	switch alert.(type) {
	case alerts.EthereumAlert:
		return theme.EthereumIcon
	case alerts.GasAlert:
		return theme.GasIcon
	case alerts.CollectionAlert:
		return theme.AlarmIcon
	}
	// This shouldn't happen
	return theme.BrokenIcon
}

type ButtonImportance int32

const (
	RegularButton ButtonImportance = iota
	PrimaryButton
	ErrorButton
)

func (theme *Theme) Button(txt string, button *widget.Clickable, importance ButtonImportance) material.ButtonStyle {
	style := material.Button(theme.Material(), button, txt)
	style.Inset = theme.MediumInset
	switch importance {
	case RegularButton:
		style.Background = theme.LighterBg // widgets.MulAlpha(theme.ContrastBg, 0.25)
	case PrimaryButton:
		// Do nothing
	case ErrorButton:
		style.Background = theme.Error
	}
	return style
}

func (theme *Theme) Hyperlink(txt string, button *widget.Clickable) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return button.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			pointer.CursorPointer.Add(gtx.Ops)
			label := material.Body1(theme.Material(), txt)
			label.Alignment = text.End
			label.Color = theme.Link
			return label.Layout(gtx)
		})
	}
}

// TODO: we should support image icons, not only vectors
func (theme *Theme) IconButton(icon *widgets.Icon, button *widget.Clickable) material.IconButtonStyle {
	wIcon := icon.WidgetIcon()
	if wIcon == nil {
		panic("can not create an icon button with image icon")
	}
	return material.IconButtonStyle{
		Background: theme.Bg,
		Color:      theme.HighImpFg,
		Icon:       wIcon,
		Size:       theme.IconSize,
		Inset:      theme.SmallInset,
		Button:     button,
	}
}

func (theme *Theme) Border(gtx layout.Context, c color.NRGBA, w layout.Widget) layout.Dimensions {
	return widgets.Border(gtx, c, theme.CornerRadius, theme.BorderWidth, w)
}

func (theme *Theme) Background(gtx layout.Context, c color.NRGBA, w layout.Widget) layout.Dimensions {
	return widgets.BackgroundRect(gtx, c, theme.CornerRadius, w)
}

/*
func (theme *Theme) Shadow(gtx layout.Context, w layout.Widget) layout.Dimensions {
	return widgets.Shadow(gtx, theme.CornerRadius, theme.ShadowElevation, w)
}
*/

func (theme *Theme) LayoutList(gtx layout.Context, list *widget.List, widgets ...layout.Widget) layout.Dimensions {
	return material.List(theme.Material(), list).Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return widgets[index](gtx)
	})
}

func (theme *Theme) LayoutListSpaced(gtx layout.Context, list *widget.List, spacer layout.Spacer, widgets ...layout.Widget) layout.Dimensions {
	count := len(widgets) + (len(widgets) - 1)
	childs := make([]layout.Widget, count)
	for i := 0; i < count; i++ {
		if i%2 == 0 {
			childs[i] = widgets[i/2]
		} else {
			childs[i] = spacer.Layout
			// func(gtx layout.Context) layout.Dimensions {
			// 	return widgets.BackgroundRect(gtx, color.NRGBA{G: 255, A: 255}, 0, spacer.Layout)
			// }
		}
	}
	return theme.LayoutList(gtx, list, childs...)
}

func (theme *Theme) LayoutForm(gtx layout.Context, list *widget.List, okButton *widget.Clickable, widgets ...layout.Widget) layout.Dimensions {
	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return theme.Button("OK", okButton, PrimaryButton).Layout(gtx)
	})
	return theme.LayoutListSpaced(gtx, list, theme.SmallVSpacer, widgets...)
}

func (theme *Theme) LayoutFlexListSpaced(gtx layout.Context, spacer layout.Spacer, widgets ...layout.FlexChild) layout.Dimensions {
	count := len(widgets) + (len(widgets) - 1)
	childs := make([]layout.FlexChild, count)
	spacerFlex := layout.Rigid(spacer.Layout)
	for i := 0; i < count; i++ {
		if i%2 == 0 {
			childs[i] = widgets[i/2]
		} else {
			childs[i] = spacerFlex
		}
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, childs...)
}

func (theme *Theme) LayoutFlexGrid(gtx layout.Context, rows, cols int, w ...layout.Widget) layout.Dimensions {
	noopWidget := func(layout.Context) layout.Dimensions { return layout.Dimensions{} }
	rowList := make([]layout.FlexChild, rows)
	for i := 0; i < rows; i++ {
		i := i
		rowList[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			colList := make([]layout.FlexChild, cols)
			for j := 0; j < cols; j++ {
				index := i*cols + j
				widget := noopWidget
				if index < len(w) {
					widget = w[index]
				}
				colList[j] = layout.Flexed(1, widget)
			}
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx, colList...)
		})
	}
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx, rowList...)
}

func rgb(c uint32) color.NRGBA {
	return rgba((c << 8) | 0xff)
}

func rgba(c uint32) color.NRGBA {
	return color.NRGBA{
		R: uint8(c >> 24),
		G: uint8(c >> 16),
		B: uint8(c >> 8),
		A: uint8(c),
	}
}

// for debugging
// TODO: these should be only included in debug build
var iconList = []struct {
	name string
	data []byte
}{
	{"Action3DRotation", icons.Action3DRotation},
	{"ActionAccessibility", icons.ActionAccessibility},
	{"ActionAccessible", icons.ActionAccessible},
	{"ActionAccountBalance", icons.ActionAccountBalance},
	{"ActionAccountBalanceWallet", icons.ActionAccountBalanceWallet},
	{"ActionAccountBox", icons.ActionAccountBox},
	{"ActionAccountCircle", icons.ActionAccountCircle},
	{"ActionAddShoppingCart", icons.ActionAddShoppingCart},
	{"ActionAlarm", icons.ActionAlarm},
	{"ActionAlarmAdd", icons.ActionAlarmAdd},
	{"ActionAlarmOff", icons.ActionAlarmOff},
	{"ActionAlarmOn", icons.ActionAlarmOn},
	{"ActionAllOut", icons.ActionAllOut},
	{"ActionAndroid", icons.ActionAndroid},
	{"ActionAnnouncement", icons.ActionAnnouncement},
	{"ActionAspectRatio", icons.ActionAspectRatio},
	{"ActionAssessment", icons.ActionAssessment},
	{"ActionAssignment", icons.ActionAssignment},
	{"ActionAssignmentInd", icons.ActionAssignmentInd},
	{"ActionAssignmentLate", icons.ActionAssignmentLate},
	{"ActionAssignmentReturn", icons.ActionAssignmentReturn},
	{"ActionAssignmentReturned", icons.ActionAssignmentReturned},
	{"ActionAssignmentTurnedIn", icons.ActionAssignmentTurnedIn},
	{"ActionAutorenew", icons.ActionAutorenew},
	{"ActionBackup", icons.ActionBackup},
	{"ActionBook", icons.ActionBook},
	{"ActionBookmark", icons.ActionBookmark},
	{"ActionBookmarkBorder", icons.ActionBookmarkBorder},
	{"ActionBugReport", icons.ActionBugReport},
	{"ActionBuild", icons.ActionBuild},
	{"ActionCached", icons.ActionCached},
	{"ActionCameraEnhance", icons.ActionCameraEnhance},
	{"ActionCardGiftcard", icons.ActionCardGiftcard},
	{"ActionCardMembership", icons.ActionCardMembership},
	{"ActionCardTravel", icons.ActionCardTravel},
	{"ActionChangeHistory", icons.ActionChangeHistory},
	{"ActionCheckCircle", icons.ActionCheckCircle},
	{"ActionChromeReaderMode", icons.ActionChromeReaderMode},
	{"ActionClass", icons.ActionClass},
	{"ActionCode", icons.ActionCode},
	{"ActionCompareArrows", icons.ActionCompareArrows},
	{"ActionCopyright", icons.ActionCopyright},
	{"ActionCreditCard", icons.ActionCreditCard},
	{"ActionDashboard", icons.ActionDashboard},
	{"ActionDateRange", icons.ActionDateRange},
	{"ActionDelete", icons.ActionDelete},
	{"ActionDeleteForever", icons.ActionDeleteForever},
	{"ActionDescription", icons.ActionDescription},
	{"ActionDNS", icons.ActionDNS},
	{"ActionDone", icons.ActionDone},
	{"ActionDoneAll", icons.ActionDoneAll},
	{"ActionDonutLarge", icons.ActionDonutLarge},
	{"ActionDonutSmall", icons.ActionDonutSmall},
	{"ActionEject", icons.ActionEject},
	{"ActionEuroSymbol", icons.ActionEuroSymbol},
	{"ActionEvent", icons.ActionEvent},
	{"ActionEventSeat", icons.ActionEventSeat},
	{"ActionExitToApp", icons.ActionExitToApp},
	{"ActionExplore", icons.ActionExplore},
	{"ActionExtension", icons.ActionExtension},
	{"ActionFace", icons.ActionFace},
	{"ActionFavorite", icons.ActionFavorite},
	{"ActionFavoriteBorder", icons.ActionFavoriteBorder},
	{"ActionFeedback", icons.ActionFeedback},
	{"ActionFindInPage", icons.ActionFindInPage},
	{"ActionFindReplace", icons.ActionFindReplace},
	{"ActionFingerprint", icons.ActionFingerprint},
	{"ActionFlightLand", icons.ActionFlightLand},
	{"ActionFlightTakeoff", icons.ActionFlightTakeoff},
	{"ActionFlipToBack", icons.ActionFlipToBack},
	{"ActionFlipToFront", icons.ActionFlipToFront},
	{"ActionGTranslate", icons.ActionGTranslate},
	{"ActionGavel", icons.ActionGavel},
	{"ActionGetApp", icons.ActionGetApp},
	{"ActionGIF", icons.ActionGIF},
	{"ActionGrade", icons.ActionGrade},
	{"ActionGroupWork", icons.ActionGroupWork},
	{"ActionHelp", icons.ActionHelp},
	{"ActionHelpOutline", icons.ActionHelpOutline},
	{"ActionHighlightOff", icons.ActionHighlightOff},
	{"ActionHistory", icons.ActionHistory},
	{"ActionHome", icons.ActionHome},
	{"ActionHourglassEmpty", icons.ActionHourglassEmpty},
	{"ActionHourglassFull", icons.ActionHourglassFull},
	{"ActionHTTP", icons.ActionHTTP},
	{"ActionHTTPS", icons.ActionHTTPS},
	{"ActionImportantDevices", icons.ActionImportantDevices},
	{"ActionInfo", icons.ActionInfo},
	{"ActionInfoOutline", icons.ActionInfoOutline},
	{"ActionInput", icons.ActionInput},
	{"ActionInvertColors", icons.ActionInvertColors},
	{"ActionLabel", icons.ActionLabel},
	{"ActionLabelOutline", icons.ActionLabelOutline},
	{"ActionLanguage", icons.ActionLanguage},
	{"ActionLaunch", icons.ActionLaunch},
	{"ActionLightbulbOutline", icons.ActionLightbulbOutline},
	{"ActionLineStyle", icons.ActionLineStyle},
	{"ActionLineWeight", icons.ActionLineWeight},
	{"ActionList", icons.ActionList},
	{"ActionLock", icons.ActionLock},
	{"ActionLockOpen", icons.ActionLockOpen},
	{"ActionLockOutline", icons.ActionLockOutline},
	{"ActionLoyalty", icons.ActionLoyalty},
	{"ActionMarkUnreadMailbox", icons.ActionMarkUnreadMailbox},
	{"ActionMotorcycle", icons.ActionMotorcycle},
	{"ActionNoteAdd", icons.ActionNoteAdd},
	{"ActionOfflinePin", icons.ActionOfflinePin},
	{"ActionOpacity", icons.ActionOpacity},
	{"ActionOpenInBrowser", icons.ActionOpenInBrowser},
	{"ActionOpenInNew", icons.ActionOpenInNew},
	{"ActionOpenWith", icons.ActionOpenWith},
	{"ActionPageview", icons.ActionPageview},
	{"ActionPanTool", icons.ActionPanTool},
	{"ActionPayment", icons.ActionPayment},
	{"ActionPermCameraMic", icons.ActionPermCameraMic},
	{"ActionPermContactCalendar", icons.ActionPermContactCalendar},
	{"ActionPermDataSetting", icons.ActionPermDataSetting},
	{"ActionPermDeviceInformation", icons.ActionPermDeviceInformation},
	{"ActionPermIdentity", icons.ActionPermIdentity},
	{"ActionPermMedia", icons.ActionPermMedia},
	{"ActionPermPhoneMsg", icons.ActionPermPhoneMsg},
	{"ActionPermScanWiFi", icons.ActionPermScanWiFi},
	{"ActionPets", icons.ActionPets},
	{"ActionPictureInPicture", icons.ActionPictureInPicture},
	{"ActionPictureInPictureAlt", icons.ActionPictureInPictureAlt},
	{"ActionPlayForWork", icons.ActionPlayForWork},
	{"ActionPolymer", icons.ActionPolymer},
	{"ActionPowerSettingsNew", icons.ActionPowerSettingsNew},
	{"ActionPregnantWoman", icons.ActionPregnantWoman},
	{"ActionPrint", icons.ActionPrint},
	{"ActionQueryBuilder", icons.ActionQueryBuilder},
	{"ActionQuestionAnswer", icons.ActionQuestionAnswer},
	{"ActionReceipt", icons.ActionReceipt},
	{"ActionRecordVoiceOver", icons.ActionRecordVoiceOver},
	{"ActionRedeem", icons.ActionRedeem},
	{"ActionRemoveShoppingCart", icons.ActionRemoveShoppingCart},
	{"ActionReorder", icons.ActionReorder},
	{"ActionReportProblem", icons.ActionReportProblem},
	{"ActionRestore", icons.ActionRestore},
	{"ActionRestorePage", icons.ActionRestorePage},
	{"ActionRoom", icons.ActionRoom},
	{"ActionRoundedCorner", icons.ActionRoundedCorner},
	{"ActionRowing", icons.ActionRowing},
	{"ActionSchedule", icons.ActionSchedule},
	{"ActionSearch", icons.ActionSearch},
	{"ActionSettings", icons.ActionSettings},
	{"ActionSettingsApplications", icons.ActionSettingsApplications},
	{"ActionSettingsBackupRestore", icons.ActionSettingsBackupRestore},
	{"ActionSettingsBluetooth", icons.ActionSettingsBluetooth},
	{"ActionSettingsBrightness", icons.ActionSettingsBrightness},
	{"ActionSettingsCell", icons.ActionSettingsCell},
	{"ActionSettingsEthernet", icons.ActionSettingsEthernet},
	{"ActionSettingsInputAntenna", icons.ActionSettingsInputAntenna},
	{"ActionSettingsInputComponent", icons.ActionSettingsInputComponent},
	{"ActionSettingsInputComposite", icons.ActionSettingsInputComposite},
	{"ActionSettingsInputHDMI", icons.ActionSettingsInputHDMI},
	{"ActionSettingsInputSVideo", icons.ActionSettingsInputSVideo},
	{"ActionSettingsOverscan", icons.ActionSettingsOverscan},
	{"ActionSettingsPhone", icons.ActionSettingsPhone},
	{"ActionSettingsPower", icons.ActionSettingsPower},
	{"ActionSettingsRemote", icons.ActionSettingsRemote},
	{"ActionSettingsVoice", icons.ActionSettingsVoice},
	{"ActionShop", icons.ActionShop},
	{"ActionShopTwo", icons.ActionShopTwo},
	{"ActionShoppingBasket", icons.ActionShoppingBasket},
	{"ActionShoppingCart", icons.ActionShoppingCart},
	{"ActionSpeakerNotes", icons.ActionSpeakerNotes},
	{"ActionSpeakerNotesOff", icons.ActionSpeakerNotesOff},
	{"ActionSpellcheck", icons.ActionSpellcheck},
	{"ActionStarRate", icons.ActionStarRate},
	{"ActionStars", icons.ActionStars},
	{"ActionStore", icons.ActionStore},
	{"ActionSubject", icons.ActionSubject},
	{"ActionSupervisorAccount", icons.ActionSupervisorAccount},
	{"ActionSwapHoriz", icons.ActionSwapHoriz},
	{"ActionSwapVert", icons.ActionSwapVert},
	{"ActionSwapVerticalCircle", icons.ActionSwapVerticalCircle},
	{"ActionSystemUpdateAlt", icons.ActionSystemUpdateAlt},
	{"ActionTab", icons.ActionTab},
	{"ActionTabUnselected", icons.ActionTabUnselected},
	{"ActionTheaters", icons.ActionTheaters},
	{"ActionThumbDown", icons.ActionThumbDown},
	{"ActionThumbUp", icons.ActionThumbUp},
	{"ActionThumbsUpDown", icons.ActionThumbsUpDown},
	{"ActionTimeline", icons.ActionTimeline},
	{"ActionTOC", icons.ActionTOC},
	{"ActionToday", icons.ActionToday},
	{"ActionToll", icons.ActionToll},
	{"ActionTouchApp", icons.ActionTouchApp},
	{"ActionTrackChanges", icons.ActionTrackChanges},
	{"ActionTranslate", icons.ActionTranslate},
	{"ActionTrendingDown", icons.ActionTrendingDown},
	{"ActionTrendingFlat", icons.ActionTrendingFlat},
	{"ActionTrendingUp", icons.ActionTrendingUp},
	{"ActionTurnedIn", icons.ActionTurnedIn},
	{"ActionTurnedInNot", icons.ActionTurnedInNot},
	{"ActionUpdate", icons.ActionUpdate},
	{"ActionVerifiedUser", icons.ActionVerifiedUser},
	{"ActionViewAgenda", icons.ActionViewAgenda},
	{"ActionViewArray", icons.ActionViewArray},
	{"ActionViewCarousel", icons.ActionViewCarousel},
	{"ActionViewColumn", icons.ActionViewColumn},
	{"ActionViewDay", icons.ActionViewDay},
	{"ActionViewHeadline", icons.ActionViewHeadline},
	{"ActionViewList", icons.ActionViewList},
	{"ActionViewModule", icons.ActionViewModule},
	{"ActionViewQuilt", icons.ActionViewQuilt},
	{"ActionViewStream", icons.ActionViewStream},
	{"ActionViewWeek", icons.ActionViewWeek},
	{"ActionVisibility", icons.ActionVisibility},
	{"ActionVisibilityOff", icons.ActionVisibilityOff},
	{"ActionWatchLater", icons.ActionWatchLater},
	{"ActionWork", icons.ActionWork},
	{"ActionYoutubeSearchedFor", icons.ActionYoutubeSearchedFor},
	{"ActionZoomIn", icons.ActionZoomIn},
	{"ActionZoomOut", icons.ActionZoomOut},
	{"AlertAddAlert", icons.AlertAddAlert},
	{"AlertError", icons.AlertError},
	{"AlertErrorOutline", icons.AlertErrorOutline},
	{"AlertWarning", icons.AlertWarning},
	{"AVAddToQueue", icons.AVAddToQueue},
	{"AVAirplay", icons.AVAirplay},
	{"AVAlbum", icons.AVAlbum},
	{"AVArtTrack", icons.AVArtTrack},
	{"AVAVTimer", icons.AVAVTimer},
	{"AVBrandingWatermark", icons.AVBrandingWatermark},
	{"AVCallToAction", icons.AVCallToAction},
	{"AVClosedCaption", icons.AVClosedCaption},
	{"AVEqualizer", icons.AVEqualizer},
	{"AVExplicit", icons.AVExplicit},
	{"AVFastForward", icons.AVFastForward},
	{"AVFastRewind", icons.AVFastRewind},
	{"AVFeaturedPlayList", icons.AVFeaturedPlayList},
	{"AVFeaturedVideo", icons.AVFeaturedVideo},
	{"AVFiberDVR", icons.AVFiberDVR},
	{"AVFiberManualRecord", icons.AVFiberManualRecord},
	{"AVFiberNew", icons.AVFiberNew},
	{"AVFiberPin", icons.AVFiberPin},
	{"AVFiberSmartRecord", icons.AVFiberSmartRecord},
	{"AVForward10", icons.AVForward10},
	{"AVForward30", icons.AVForward30},
	{"AVForward5", icons.AVForward5},
	{"AVGames", icons.AVGames},
	{"AVHD", icons.AVHD},
	{"AVHearing", icons.AVHearing},
	{"AVHighQuality", icons.AVHighQuality},
	{"AVLibraryAdd", icons.AVLibraryAdd},
	{"AVLibraryBooks", icons.AVLibraryBooks},
	{"AVLibraryMusic", icons.AVLibraryMusic},
	{"AVLoop", icons.AVLoop},
	{"AVMic", icons.AVMic},
	{"AVMicNone", icons.AVMicNone},
	{"AVMicOff", icons.AVMicOff},
	{"AVMovie", icons.AVMovie},
	{"AVMusicVideo", icons.AVMusicVideo},
	{"AVNewReleases", icons.AVNewReleases},
	{"AVNotInterested", icons.AVNotInterested},
	{"AVNote", icons.AVNote},
	{"AVPause", icons.AVPause},
	{"AVPauseCircleFilled", icons.AVPauseCircleFilled},
	{"AVPauseCircleOutline", icons.AVPauseCircleOutline},
	{"AVPlayArrow", icons.AVPlayArrow},
	{"AVPlayCircleFilled", icons.AVPlayCircleFilled},
	{"AVPlayCircleOutline", icons.AVPlayCircleOutline},
	{"AVPlaylistAdd", icons.AVPlaylistAdd},
	{"AVPlaylistAddCheck", icons.AVPlaylistAddCheck},
	{"AVPlaylistPlay", icons.AVPlaylistPlay},
	{"AVQueue", icons.AVQueue},
	{"AVQueueMusic", icons.AVQueueMusic},
	{"AVQueuePlayNext", icons.AVQueuePlayNext},
	{"AVRadio", icons.AVRadio},
	{"AVRecentActors", icons.AVRecentActors},
	{"AVRemoveFromQueue", icons.AVRemoveFromQueue},
	{"AVRepeat", icons.AVRepeat},
	{"AVRepeatOne", icons.AVRepeatOne},
	{"AVReplay", icons.AVReplay},
	{"AVReplay10", icons.AVReplay10},
	{"AVReplay30", icons.AVReplay30},
	{"AVReplay5", icons.AVReplay5},
	{"AVShuffle", icons.AVShuffle},
	{"AVSkipNext", icons.AVSkipNext},
	{"AVSkipPrevious", icons.AVSkipPrevious},
	{"AVSlowMotionVideo", icons.AVSlowMotionVideo},
	{"AVSnooze", icons.AVSnooze},
	{"AVSortByAlpha", icons.AVSortByAlpha},
	{"AVStop", icons.AVStop},
	{"AVSubscriptions", icons.AVSubscriptions},
	{"AVSubtitles", icons.AVSubtitles},
	{"AVSurroundSound", icons.AVSurroundSound},
	{"AVVideoCall", icons.AVVideoCall},
	{"AVVideoLabel", icons.AVVideoLabel},
	{"AVVideoLibrary", icons.AVVideoLibrary},
	{"AVVideocam", icons.AVVideocam},
	{"AVVideocamOff", icons.AVVideocamOff},
	{"AVVolumeDown", icons.AVVolumeDown},
	{"AVVolumeMute", icons.AVVolumeMute},
	{"AVVolumeOff", icons.AVVolumeOff},
	{"AVVolumeUp", icons.AVVolumeUp},
	{"AVWeb", icons.AVWeb},
	{"AVWebAsset", icons.AVWebAsset},
	{"CommunicationBusiness", icons.CommunicationBusiness},
	{"CommunicationCall", icons.CommunicationCall},
	{"CommunicationCallEnd", icons.CommunicationCallEnd},
	{"CommunicationCallMade", icons.CommunicationCallMade},
	{"CommunicationCallMerge", icons.CommunicationCallMerge},
	{"CommunicationCallMissed", icons.CommunicationCallMissed},
	{"CommunicationCallMissedOutgoing", icons.CommunicationCallMissedOutgoing},
	{"CommunicationCallReceived", icons.CommunicationCallReceived},
	{"CommunicationCallSplit", icons.CommunicationCallSplit},
	{"CommunicationChat", icons.CommunicationChat},
	{"CommunicationChatBubble", icons.CommunicationChatBubble},
	{"CommunicationChatBubbleOutline", icons.CommunicationChatBubbleOutline},
	{"CommunicationClearAll", icons.CommunicationClearAll},
	{"CommunicationComment", icons.CommunicationComment},
	{"CommunicationContactMail", icons.CommunicationContactMail},
	{"CommunicationContactPhone", icons.CommunicationContactPhone},
	{"CommunicationContacts", icons.CommunicationContacts},
	{"CommunicationDialerSIP", icons.CommunicationDialerSIP},
	{"CommunicationDialpad", icons.CommunicationDialpad},
	{"CommunicationEmail", icons.CommunicationEmail},
	{"CommunicationForum", icons.CommunicationForum},
	{"CommunicationImportContacts", icons.CommunicationImportContacts},
	{"CommunicationImportExport", icons.CommunicationImportExport},
	{"CommunicationInvertColorsOff", icons.CommunicationInvertColorsOff},
	{"CommunicationLiveHelp", icons.CommunicationLiveHelp},
	{"CommunicationLocationOff", icons.CommunicationLocationOff},
	{"CommunicationLocationOn", icons.CommunicationLocationOn},
	{"CommunicationMailOutline", icons.CommunicationMailOutline},
	{"CommunicationMessage", icons.CommunicationMessage},
	{"CommunicationNoSIM", icons.CommunicationNoSIM},
	{"CommunicationPhone", icons.CommunicationPhone},
	{"CommunicationPhoneLinkErase", icons.CommunicationPhoneLinkErase},
	{"CommunicationPhoneLinkLock", icons.CommunicationPhoneLinkLock},
	{"CommunicationPhoneLinkRing", icons.CommunicationPhoneLinkRing},
	{"CommunicationPhoneLinkSetup", icons.CommunicationPhoneLinkSetup},
	{"CommunicationPortableWiFiOff", icons.CommunicationPortableWiFiOff},
	{"CommunicationPresentToAll", icons.CommunicationPresentToAll},
	{"CommunicationRingVolume", icons.CommunicationRingVolume},
	{"CommunicationRSSFeed", icons.CommunicationRSSFeed},
	{"CommunicationScreenShare", icons.CommunicationScreenShare},
	{"CommunicationSpeakerPhone", icons.CommunicationSpeakerPhone},
	{"CommunicationStayCurrentLandscape", icons.CommunicationStayCurrentLandscape},
	{"CommunicationStayCurrentPortrait", icons.CommunicationStayCurrentPortrait},
	{"CommunicationStayPrimaryLandscape", icons.CommunicationStayPrimaryLandscape},
	{"CommunicationStayPrimaryPortrait", icons.CommunicationStayPrimaryPortrait},
	{"CommunicationStopScreenShare", icons.CommunicationStopScreenShare},
	{"CommunicationSwapCalls", icons.CommunicationSwapCalls},
	{"CommunicationTextSMS", icons.CommunicationTextSMS},
	{"CommunicationVoicemail", icons.CommunicationVoicemail},
	{"CommunicationVPNKey", icons.CommunicationVPNKey},
	{"ContentAdd", icons.ContentAdd},
	{"ContentAddBox", icons.ContentAddBox},
	{"ContentAddCircle", icons.ContentAddCircle},
	{"ContentAddCircleOutline", icons.ContentAddCircleOutline},
	{"ContentArchive", icons.ContentArchive},
	{"ContentBackspace", icons.ContentBackspace},
	{"ContentBlock", icons.ContentBlock},
	{"ContentClear", icons.ContentClear},
	{"ContentContentCopy", icons.ContentContentCopy},
	{"ContentContentCut", icons.ContentContentCut},
	{"ContentContentPaste", icons.ContentContentPaste},
	{"ContentCreate", icons.ContentCreate},
	{"ContentDeleteSweep", icons.ContentDeleteSweep},
	{"ContentDrafts", icons.ContentDrafts},
	{"ContentFilterList", icons.ContentFilterList},
	{"ContentFlag", icons.ContentFlag},
	{"ContentFontDownload", icons.ContentFontDownload},
	{"ContentForward", icons.ContentForward},
	{"ContentGesture", icons.ContentGesture},
	{"ContentInbox", icons.ContentInbox},
	{"ContentLink", icons.ContentLink},
	{"ContentLowPriority", icons.ContentLowPriority},
	{"ContentMail", icons.ContentMail},
	{"ContentMarkUnread", icons.ContentMarkUnread},
	{"ContentMoveToInbox", icons.ContentMoveToInbox},
	{"ContentNextWeek", icons.ContentNextWeek},
	{"ContentRedo", icons.ContentRedo},
	{"ContentRemove", icons.ContentRemove},
	{"ContentRemoveCircle", icons.ContentRemoveCircle},
	{"ContentRemoveCircleOutline", icons.ContentRemoveCircleOutline},
	{"ContentReply", icons.ContentReply},
	{"ContentReplyAll", icons.ContentReplyAll},
	{"ContentReport", icons.ContentReport},
	{"ContentSave", icons.ContentSave},
	{"ContentSelectAll", icons.ContentSelectAll},
	{"ContentSend", icons.ContentSend},
	{"ContentSort", icons.ContentSort},
	{"ContentTextFormat", icons.ContentTextFormat},
	{"ContentUnarchive", icons.ContentUnarchive},
	{"ContentUndo", icons.ContentUndo},
	{"ContentWeekend", icons.ContentWeekend},
	{"DeviceAccessAlarm", icons.DeviceAccessAlarm},
	{"DeviceAccessAlarms", icons.DeviceAccessAlarms},
	{"DeviceAccessTime", icons.DeviceAccessTime},
	{"DeviceAddAlarm", icons.DeviceAddAlarm},
	{"DeviceAirplaneModeActive", icons.DeviceAirplaneModeActive},
	{"DeviceAirplaneModeInactive", icons.DeviceAirplaneModeInactive},
	{"DeviceBattery20", icons.DeviceBattery20},
	{"DeviceBattery30", icons.DeviceBattery30},
	{"DeviceBattery50", icons.DeviceBattery50},
	{"DeviceBattery60", icons.DeviceBattery60},
	{"DeviceBattery80", icons.DeviceBattery80},
	{"DeviceBattery90", icons.DeviceBattery90},
	{"DeviceBatteryAlert", icons.DeviceBatteryAlert},
	{"DeviceBatteryCharging20", icons.DeviceBatteryCharging20},
	{"DeviceBatteryCharging30", icons.DeviceBatteryCharging30},
	{"DeviceBatteryCharging50", icons.DeviceBatteryCharging50},
	{"DeviceBatteryCharging60", icons.DeviceBatteryCharging60},
	{"DeviceBatteryCharging80", icons.DeviceBatteryCharging80},
	{"DeviceBatteryCharging90", icons.DeviceBatteryCharging90},
	{"DeviceBatteryChargingFull", icons.DeviceBatteryChargingFull},
	{"DeviceBatteryFull", icons.DeviceBatteryFull},
	{"DeviceBatteryStd", icons.DeviceBatteryStd},
	{"DeviceBatteryUnknown", icons.DeviceBatteryUnknown},
	{"DeviceBluetooth", icons.DeviceBluetooth},
	{"DeviceBluetoothConnected", icons.DeviceBluetoothConnected},
	{"DeviceBluetoothDisabled", icons.DeviceBluetoothDisabled},
	{"DeviceBluetoothSearching", icons.DeviceBluetoothSearching},
	{"DeviceBrightnessAuto", icons.DeviceBrightnessAuto},
	{"DeviceBrightnessHigh", icons.DeviceBrightnessHigh},
	{"DeviceBrightnessLow", icons.DeviceBrightnessLow},
	{"DeviceBrightnessMedium", icons.DeviceBrightnessMedium},
	{"DeviceDataUsage", icons.DeviceDataUsage},
	{"DeviceDeveloperMode", icons.DeviceDeveloperMode},
	{"DeviceDevices", icons.DeviceDevices},
	{"DeviceDVR", icons.DeviceDVR},
	{"DeviceGPSFixed", icons.DeviceGPSFixed},
	{"DeviceGPSNotFixed", icons.DeviceGPSNotFixed},
	{"DeviceGPSOff", icons.DeviceGPSOff},
	{"DeviceGraphicEq", icons.DeviceGraphicEq},
	{"DeviceLocationDisabled", icons.DeviceLocationDisabled},
	{"DeviceLocationSearching", icons.DeviceLocationSearching},
	{"DeviceNetworkCell", icons.DeviceNetworkCell},
	{"DeviceNetworkWiFi", icons.DeviceNetworkWiFi},
	{"DeviceNFC", icons.DeviceNFC},
	{"DeviceScreenLockLandscape", icons.DeviceScreenLockLandscape},
	{"DeviceScreenLockPortrait", icons.DeviceScreenLockPortrait},
	{"DeviceScreenLockRotation", icons.DeviceScreenLockRotation},
	{"DeviceScreenRotation", icons.DeviceScreenRotation},
	{"DeviceSDStorage", icons.DeviceSDStorage},
	{"DeviceSettingsSystemDaydream", icons.DeviceSettingsSystemDaydream},
	{"DeviceSignalCellular0Bar", icons.DeviceSignalCellular0Bar},
	{"DeviceSignalCellular1Bar", icons.DeviceSignalCellular1Bar},
	{"DeviceSignalCellular2Bar", icons.DeviceSignalCellular2Bar},
	{"DeviceSignalCellular3Bar", icons.DeviceSignalCellular3Bar},
	{"DeviceSignalCellular4Bar", icons.DeviceSignalCellular4Bar},
	{"DeviceSignalCellularConnectedNoInternet0Bar", icons.DeviceSignalCellularConnectedNoInternet0Bar},
	{"DeviceSignalCellularConnectedNoInternet1Bar", icons.DeviceSignalCellularConnectedNoInternet1Bar},
	{"DeviceSignalCellularConnectedNoInternet2Bar", icons.DeviceSignalCellularConnectedNoInternet2Bar},
	{"DeviceSignalCellularConnectedNoInternet3Bar", icons.DeviceSignalCellularConnectedNoInternet3Bar},
	{"DeviceSignalCellularConnectedNoInternet4Bar", icons.DeviceSignalCellularConnectedNoInternet4Bar},
	{"DeviceSignalCellularNoSIM", icons.DeviceSignalCellularNoSIM},
	{"DeviceSignalCellularNull", icons.DeviceSignalCellularNull},
	{"DeviceSignalCellularOff", icons.DeviceSignalCellularOff},
	{"DeviceSignalWiFi0Bar", icons.DeviceSignalWiFi0Bar},
	{"DeviceSignalWiFi1Bar", icons.DeviceSignalWiFi1Bar},
	{"DeviceSignalWiFi1BarLock", icons.DeviceSignalWiFi1BarLock},
	{"DeviceSignalWiFi2Bar", icons.DeviceSignalWiFi2Bar},
	{"DeviceSignalWiFi2BarLock", icons.DeviceSignalWiFi2BarLock},
	{"DeviceSignalWiFi3Bar", icons.DeviceSignalWiFi3Bar},
	{"DeviceSignalWiFi3BarLock", icons.DeviceSignalWiFi3BarLock},
	{"DeviceSignalWiFi4Bar", icons.DeviceSignalWiFi4Bar},
	{"DeviceSignalWiFi4BarLock", icons.DeviceSignalWiFi4BarLock},
	{"DeviceSignalWiFiOff", icons.DeviceSignalWiFiOff},
	{"DeviceStorage", icons.DeviceStorage},
	{"DeviceUSB", icons.DeviceUSB},
	{"DeviceWallpaper", icons.DeviceWallpaper},
	{"DeviceWidgets", icons.DeviceWidgets},
	{"DeviceWiFiLock", icons.DeviceWiFiLock},
	{"DeviceWiFiTethering", icons.DeviceWiFiTethering},
	{"EditorAttachFile", icons.EditorAttachFile},
	{"EditorAttachMoney", icons.EditorAttachMoney},
	{"EditorBorderAll", icons.EditorBorderAll},
	{"EditorBorderBottom", icons.EditorBorderBottom},
	{"EditorBorderClear", icons.EditorBorderClear},
	{"EditorBorderColor", icons.EditorBorderColor},
	{"EditorBorderHorizontal", icons.EditorBorderHorizontal},
	{"EditorBorderInner", icons.EditorBorderInner},
	{"EditorBorderLeft", icons.EditorBorderLeft},
	{"EditorBorderOuter", icons.EditorBorderOuter},
	{"EditorBorderRight", icons.EditorBorderRight},
	{"EditorBorderStyle", icons.EditorBorderStyle},
	{"EditorBorderTop", icons.EditorBorderTop},
	{"EditorBorderVertical", icons.EditorBorderVertical},
	{"EditorBubbleChart", icons.EditorBubbleChart},
	{"EditorDragHandle", icons.EditorDragHandle},
	{"EditorFormatAlignCenter", icons.EditorFormatAlignCenter},
	{"EditorFormatAlignJustify", icons.EditorFormatAlignJustify},
	{"EditorFormatAlignLeft", icons.EditorFormatAlignLeft},
	{"EditorFormatAlignRight", icons.EditorFormatAlignRight},
	{"EditorFormatBold", icons.EditorFormatBold},
	{"EditorFormatClear", icons.EditorFormatClear},
	{"EditorFormatColorFill", icons.EditorFormatColorFill},
	{"EditorFormatColorReset", icons.EditorFormatColorReset},
	{"EditorFormatColorText", icons.EditorFormatColorText},
	{"EditorFormatIndentDecrease", icons.EditorFormatIndentDecrease},
	{"EditorFormatIndentIncrease", icons.EditorFormatIndentIncrease},
	{"EditorFormatItalic", icons.EditorFormatItalic},
	{"EditorFormatLineSpacing", icons.EditorFormatLineSpacing},
	{"EditorFormatListBulleted", icons.EditorFormatListBulleted},
	{"EditorFormatListNumbered", icons.EditorFormatListNumbered},
	{"EditorFormatPaint", icons.EditorFormatPaint},
	{"EditorFormatQuote", icons.EditorFormatQuote},
	{"EditorFormatShapes", icons.EditorFormatShapes},
	{"EditorFormatSize", icons.EditorFormatSize},
	{"EditorFormatStrikethrough", icons.EditorFormatStrikethrough},
	{"EditorFormatTextDirectionLToR", icons.EditorFormatTextDirectionLToR},
	{"EditorFormatTextDirectionRToL", icons.EditorFormatTextDirectionRToL},
	{"EditorFormatUnderlined", icons.EditorFormatUnderlined},
	{"EditorFunctions", icons.EditorFunctions},
	{"EditorHighlight", icons.EditorHighlight},
	{"EditorInsertChart", icons.EditorInsertChart},
	{"EditorInsertComment", icons.EditorInsertComment},
	{"EditorInsertDriveFile", icons.EditorInsertDriveFile},
	{"EditorInsertEmoticon", icons.EditorInsertEmoticon},
	{"EditorInsertInvitation", icons.EditorInsertInvitation},
	{"EditorInsertLink", icons.EditorInsertLink},
	{"EditorInsertPhoto", icons.EditorInsertPhoto},
	{"EditorLinearScale", icons.EditorLinearScale},
	{"EditorMergeType", icons.EditorMergeType},
	{"EditorModeComment", icons.EditorModeComment},
	{"EditorModeEdit", icons.EditorModeEdit},
	{"EditorMonetizationOn", icons.EditorMonetizationOn},
	{"EditorMoneyOff", icons.EditorMoneyOff},
	{"EditorMultilineChart", icons.EditorMultilineChart},
	{"EditorPieChart", icons.EditorPieChart},
	{"EditorPieChartOutlined", icons.EditorPieChartOutlined},
	{"EditorPublish", icons.EditorPublish},
	{"EditorShortText", icons.EditorShortText},
	{"EditorShowChart", icons.EditorShowChart},
	{"EditorSpaceBar", icons.EditorSpaceBar},
	{"EditorStrikethroughS", icons.EditorStrikethroughS},
	{"EditorTextFields", icons.EditorTextFields},
	{"EditorTitle", icons.EditorTitle},
	{"EditorVerticalAlignBottom", icons.EditorVerticalAlignBottom},
	{"EditorVerticalAlignCenter", icons.EditorVerticalAlignCenter},
	{"EditorVerticalAlignTop", icons.EditorVerticalAlignTop},
	{"EditorWrapText", icons.EditorWrapText},
	{"FileAttachment", icons.FileAttachment},
	{"FileCloud", icons.FileCloud},
	{"FileCloudCircle", icons.FileCloudCircle},
	{"FileCloudDone", icons.FileCloudDone},
	{"FileCloudDownload", icons.FileCloudDownload},
	{"FileCloudOff", icons.FileCloudOff},
	{"FileCloudQueue", icons.FileCloudQueue},
	{"FileCloudUpload", icons.FileCloudUpload},
	{"FileCreateNewFolder", icons.FileCreateNewFolder},
	{"FileFileDownload", icons.FileFileDownload},
	{"FileFileUpload", icons.FileFileUpload},
	{"FileFolder", icons.FileFolder},
	{"FileFolderOpen", icons.FileFolderOpen},
	{"FileFolderShared", icons.FileFolderShared},
	{"HardwareCast", icons.HardwareCast},
	{"HardwareCastConnected", icons.HardwareCastConnected},
	{"HardwareComputer", icons.HardwareComputer},
	{"HardwareDesktopMac", icons.HardwareDesktopMac},
	{"HardwareDesktopWindows", icons.HardwareDesktopWindows},
	{"HardwareDeveloperBoard", icons.HardwareDeveloperBoard},
	{"HardwareDeviceHub", icons.HardwareDeviceHub},
	{"HardwareDevicesOther", icons.HardwareDevicesOther},
	{"HardwareDock", icons.HardwareDock},
	{"HardwareGamepad", icons.HardwareGamepad},
	{"HardwareHeadset", icons.HardwareHeadset},
	{"HardwareHeadsetMic", icons.HardwareHeadsetMic},
	{"HardwareKeyboard", icons.HardwareKeyboard},
	{"HardwareKeyboardArrowDown", icons.HardwareKeyboardArrowDown},
	{"HardwareKeyboardArrowLeft", icons.HardwareKeyboardArrowLeft},
	{"HardwareKeyboardArrowRight", icons.HardwareKeyboardArrowRight},
	{"HardwareKeyboardArrowUp", icons.HardwareKeyboardArrowUp},
	{"HardwareKeyboardBackspace", icons.HardwareKeyboardBackspace},
	{"HardwareKeyboardCapslock", icons.HardwareKeyboardCapslock},
	{"HardwareKeyboardHide", icons.HardwareKeyboardHide},
	{"HardwareKeyboardReturn", icons.HardwareKeyboardReturn},
	{"HardwareKeyboardTab", icons.HardwareKeyboardTab},
	{"HardwareKeyboardVoice", icons.HardwareKeyboardVoice},
	{"HardwareLaptop", icons.HardwareLaptop},
	{"HardwareLaptopChromebook", icons.HardwareLaptopChromebook},
	{"HardwareLaptopMac", icons.HardwareLaptopMac},
	{"HardwareLaptopWindows", icons.HardwareLaptopWindows},
	{"HardwareMemory", icons.HardwareMemory},
	{"HardwareMouse", icons.HardwareMouse},
	{"HardwarePhoneAndroid", icons.HardwarePhoneAndroid},
	{"HardwarePhoneIPhone", icons.HardwarePhoneIPhone},
	{"HardwarePhoneLink", icons.HardwarePhoneLink},
	{"HardwarePhoneLinkOff", icons.HardwarePhoneLinkOff},
	{"HardwarePowerInput", icons.HardwarePowerInput},
	{"HardwareRouter", icons.HardwareRouter},
	{"HardwareScanner", icons.HardwareScanner},
	{"HardwareSecurity", icons.HardwareSecurity},
	{"HardwareSIMCard", icons.HardwareSIMCard},
	{"HardwareSmartphone", icons.HardwareSmartphone},
	{"HardwareSpeaker", icons.HardwareSpeaker},
	{"HardwareSpeakerGroup", icons.HardwareSpeakerGroup},
	{"HardwareTablet", icons.HardwareTablet},
	{"HardwareTabletAndroid", icons.HardwareTabletAndroid},
	{"HardwareTabletMac", icons.HardwareTabletMac},
	{"HardwareToys", icons.HardwareToys},
	{"HardwareTV", icons.HardwareTV},
	{"HardwareVideogameAsset", icons.HardwareVideogameAsset},
	{"HardwareWatch", icons.HardwareWatch},
	{"ImageAddAPhoto", icons.ImageAddAPhoto},
	{"ImageAddToPhotos", icons.ImageAddToPhotos},
	{"ImageAdjust", icons.ImageAdjust},
	{"ImageAssistant", icons.ImageAssistant},
	{"ImageAssistantPhoto", icons.ImageAssistantPhoto},
	{"ImageAudiotrack", icons.ImageAudiotrack},
	{"ImageBlurCircular", icons.ImageBlurCircular},
	{"ImageBlurLinear", icons.ImageBlurLinear},
	{"ImageBlurOff", icons.ImageBlurOff},
	{"ImageBlurOn", icons.ImageBlurOn},
	{"ImageBrightness1", icons.ImageBrightness1},
	{"ImageBrightness2", icons.ImageBrightness2},
	{"ImageBrightness3", icons.ImageBrightness3},
	{"ImageBrightness4", icons.ImageBrightness4},
	{"ImageBrightness5", icons.ImageBrightness5},
	{"ImageBrightness6", icons.ImageBrightness6},
	{"ImageBrightness7", icons.ImageBrightness7},
	{"ImageBrokenImage", icons.ImageBrokenImage},
	{"ImageBrush", icons.ImageBrush},
	{"ImageBurstMode", icons.ImageBurstMode},
	{"ImageCamera", icons.ImageCamera},
	{"ImageCameraAlt", icons.ImageCameraAlt},
	{"ImageCameraFront", icons.ImageCameraFront},
	{"ImageCameraRear", icons.ImageCameraRear},
	{"ImageCameraRoll", icons.ImageCameraRoll},
	{"ImageCenterFocusStrong", icons.ImageCenterFocusStrong},
	{"ImageCenterFocusWeak", icons.ImageCenterFocusWeak},
	{"ImageCollections", icons.ImageCollections},
	{"ImageCollectionsBookmark", icons.ImageCollectionsBookmark},
	{"ImageColorLens", icons.ImageColorLens},
	{"ImageColorize", icons.ImageColorize},
	{"ImageCompare", icons.ImageCompare},
	{"ImageControlPoint", icons.ImageControlPoint},
	{"ImageControlPointDuplicate", icons.ImageControlPointDuplicate},
	{"ImageCrop", icons.ImageCrop},
	{"ImageCrop169", icons.ImageCrop169},
	{"ImageCrop32", icons.ImageCrop32},
	{"ImageCrop54", icons.ImageCrop54},
	{"ImageCrop75", icons.ImageCrop75},
	{"ImageCropDIN", icons.ImageCropDIN},
	{"ImageCropFree", icons.ImageCropFree},
	{"ImageCropLandscape", icons.ImageCropLandscape},
	{"ImageCropOriginal", icons.ImageCropOriginal},
	{"ImageCropPortrait", icons.ImageCropPortrait},
	{"ImageCropRotate", icons.ImageCropRotate},
	{"ImageCropSquare", icons.ImageCropSquare},
	{"ImageDehaze", icons.ImageDehaze},
	{"ImageDetails", icons.ImageDetails},
	{"ImageEdit", icons.ImageEdit},
	{"ImageExposure", icons.ImageExposure},
	{"ImageExposureNeg1", icons.ImageExposureNeg1},
	{"ImageExposureNeg2", icons.ImageExposureNeg2},
	{"ImageExposurePlus1", icons.ImageExposurePlus1},
	{"ImageExposurePlus2", icons.ImageExposurePlus2},
	{"ImageExposureZero", icons.ImageExposureZero},
	{"ImageFilter", icons.ImageFilter},
	{"ImageFilter1", icons.ImageFilter1},
	{"ImageFilter2", icons.ImageFilter2},
	{"ImageFilter3", icons.ImageFilter3},
	{"ImageFilter4", icons.ImageFilter4},
	{"ImageFilter5", icons.ImageFilter5},
	{"ImageFilter6", icons.ImageFilter6},
	{"ImageFilter7", icons.ImageFilter7},
	{"ImageFilter8", icons.ImageFilter8},
	{"ImageFilter9", icons.ImageFilter9},
	{"ImageFilter9Plus", icons.ImageFilter9Plus},
	{"ImageFilterBAndW", icons.ImageFilterBAndW},
	{"ImageFilterCenterFocus", icons.ImageFilterCenterFocus},
	{"ImageFilterDrama", icons.ImageFilterDrama},
	{"ImageFilterFrames", icons.ImageFilterFrames},
	{"ImageFilterHDR", icons.ImageFilterHDR},
	{"ImageFilterNone", icons.ImageFilterNone},
	{"ImageFilterTiltShift", icons.ImageFilterTiltShift},
	{"ImageFilterVintage", icons.ImageFilterVintage},
	{"ImageFlare", icons.ImageFlare},
	{"ImageFlashAuto", icons.ImageFlashAuto},
	{"ImageFlashOff", icons.ImageFlashOff},
	{"ImageFlashOn", icons.ImageFlashOn},
	{"ImageFlip", icons.ImageFlip},
	{"ImageGradient", icons.ImageGradient},
	{"ImageGrain", icons.ImageGrain},
	{"ImageGridOff", icons.ImageGridOff},
	{"ImageGridOn", icons.ImageGridOn},
	{"ImageHDROff", icons.ImageHDROff},
	{"ImageHDROn", icons.ImageHDROn},
	{"ImageHDRStrong", icons.ImageHDRStrong},
	{"ImageHDRWeak", icons.ImageHDRWeak},
	{"ImageHealing", icons.ImageHealing},
	{"ImageImage", icons.ImageImage},
	{"ImageImageAspectRatio", icons.ImageImageAspectRatio},
	{"ImageISO", icons.ImageISO},
	{"ImageLandscape", icons.ImageLandscape},
	{"ImageLeakAdd", icons.ImageLeakAdd},
	{"ImageLeakRemove", icons.ImageLeakRemove},
	{"ImageLens", icons.ImageLens},
	{"ImageLinkedCamera", icons.ImageLinkedCamera},
	{"ImageLooks", icons.ImageLooks},
	{"ImageLooks3", icons.ImageLooks3},
	{"ImageLooks4", icons.ImageLooks4},
	{"ImageLooks5", icons.ImageLooks5},
	{"ImageLooks6", icons.ImageLooks6},
	{"ImageLooksOne", icons.ImageLooksOne},
	{"ImageLooksTwo", icons.ImageLooksTwo},
	{"ImageLoupe", icons.ImageLoupe},
	{"ImageMonochromePhotos", icons.ImageMonochromePhotos},
	{"ImageMovieCreation", icons.ImageMovieCreation},
	{"ImageMovieFilter", icons.ImageMovieFilter},
	{"ImageMusicNote", icons.ImageMusicNote},
	{"ImageNature", icons.ImageNature},
	{"ImageNaturePeople", icons.ImageNaturePeople},
	{"ImageNavigateBefore", icons.ImageNavigateBefore},
	{"ImageNavigateNext", icons.ImageNavigateNext},
	{"ImagePalette", icons.ImagePalette},
	{"ImagePanorama", icons.ImagePanorama},
	{"ImagePanoramaFishEye", icons.ImagePanoramaFishEye},
	{"ImagePanoramaHorizontal", icons.ImagePanoramaHorizontal},
	{"ImagePanoramaVertical", icons.ImagePanoramaVertical},
	{"ImagePanoramaWideAngle", icons.ImagePanoramaWideAngle},
	{"ImagePhoto", icons.ImagePhoto},
	{"ImagePhotoAlbum", icons.ImagePhotoAlbum},
	{"ImagePhotoCamera", icons.ImagePhotoCamera},
	{"ImagePhotoFilter", icons.ImagePhotoFilter},
	{"ImagePhotoLibrary", icons.ImagePhotoLibrary},
	{"ImagePhotoSizeSelectActual", icons.ImagePhotoSizeSelectActual},
	{"ImagePhotoSizeSelectLarge", icons.ImagePhotoSizeSelectLarge},
	{"ImagePhotoSizeSelectSmall", icons.ImagePhotoSizeSelectSmall},
	{"ImagePictureAsPDF", icons.ImagePictureAsPDF},
	{"ImagePortrait", icons.ImagePortrait},
	{"ImageRemoveRedEye", icons.ImageRemoveRedEye},
	{"ImageRotate90DegreesCCW", icons.ImageRotate90DegreesCCW},
	{"ImageRotateLeft", icons.ImageRotateLeft},
	{"ImageRotateRight", icons.ImageRotateRight},
	{"ImageSlideshow", icons.ImageSlideshow},
	{"ImageStraighten", icons.ImageStraighten},
	{"ImageStyle", icons.ImageStyle},
	{"ImageSwitchCamera", icons.ImageSwitchCamera},
	{"ImageSwitchVideo", icons.ImageSwitchVideo},
	{"ImageTagFaces", icons.ImageTagFaces},
	{"ImageTexture", icons.ImageTexture},
	{"ImageTimeLapse", icons.ImageTimeLapse},
	{"ImageTimer", icons.ImageTimer},
	{"ImageTimer10", icons.ImageTimer10},
	{"ImageTimer3", icons.ImageTimer3},
	{"ImageTimerOff", icons.ImageTimerOff},
	{"ImageTonality", icons.ImageTonality},
	{"ImageTransform", icons.ImageTransform},
	{"ImageTune", icons.ImageTune},
	{"ImageViewComfy", icons.ImageViewComfy},
	{"ImageViewCompact", icons.ImageViewCompact},
	{"ImageVignette", icons.ImageVignette},
	{"ImageWBAuto", icons.ImageWBAuto},
	{"ImageWBCloudy", icons.ImageWBCloudy},
	{"ImageWBIncandescent", icons.ImageWBIncandescent},
	{"ImageWBIridescent", icons.ImageWBIridescent},
	{"ImageWBSunny", icons.ImageWBSunny},
	{"MapsAddLocation", icons.MapsAddLocation},
	{"MapsBeenhere", icons.MapsBeenhere},
	{"MapsDirections", icons.MapsDirections},
	{"MapsDirectionsBike", icons.MapsDirectionsBike},
	{"MapsDirectionsBoat", icons.MapsDirectionsBoat},
	{"MapsDirectionsBus", icons.MapsDirectionsBus},
	{"MapsDirectionsCar", icons.MapsDirectionsCar},
	{"MapsDirectionsRailway", icons.MapsDirectionsRailway},
	{"MapsDirectionsRun", icons.MapsDirectionsRun},
	{"MapsDirectionsSubway", icons.MapsDirectionsSubway},
	{"MapsDirectionsTransit", icons.MapsDirectionsTransit},
	{"MapsDirectionsWalk", icons.MapsDirectionsWalk},
	{"MapsEditLocation", icons.MapsEditLocation},
	{"MapsEVStation", icons.MapsEVStation},
	{"MapsFlight", icons.MapsFlight},
	{"MapsHotel", icons.MapsHotel},
	{"MapsLayers", icons.MapsLayers},
	{"MapsLayersClear", icons.MapsLayersClear},
	{"MapsLocalActivity", icons.MapsLocalActivity},
	{"MapsLocalAirport", icons.MapsLocalAirport},
	{"MapsLocalATM", icons.MapsLocalATM},
	{"MapsLocalBar", icons.MapsLocalBar},
	{"MapsLocalCafe", icons.MapsLocalCafe},
	{"MapsLocalCarWash", icons.MapsLocalCarWash},
	{"MapsLocalConvenienceStore", icons.MapsLocalConvenienceStore},
	{"MapsLocalDining", icons.MapsLocalDining},
	{"MapsLocalDrink", icons.MapsLocalDrink},
	{"MapsLocalFlorist", icons.MapsLocalFlorist},
	{"MapsLocalGasStation", icons.MapsLocalGasStation},
	{"MapsLocalGroceryStore", icons.MapsLocalGroceryStore},
	{"MapsLocalHospital", icons.MapsLocalHospital},
	{"MapsLocalHotel", icons.MapsLocalHotel},
	{"MapsLocalLaundryService", icons.MapsLocalLaundryService},
	{"MapsLocalLibrary", icons.MapsLocalLibrary},
	{"MapsLocalMall", icons.MapsLocalMall},
	{"MapsLocalMovies", icons.MapsLocalMovies},
	{"MapsLocalOffer", icons.MapsLocalOffer},
	{"MapsLocalParking", icons.MapsLocalParking},
	{"MapsLocalPharmacy", icons.MapsLocalPharmacy},
	{"MapsLocalPhone", icons.MapsLocalPhone},
	{"MapsLocalPizza", icons.MapsLocalPizza},
	{"MapsLocalPlay", icons.MapsLocalPlay},
	{"MapsLocalPostOffice", icons.MapsLocalPostOffice},
	{"MapsLocalPrintshop", icons.MapsLocalPrintshop},
	{"MapsLocalSee", icons.MapsLocalSee},
	{"MapsLocalShipping", icons.MapsLocalShipping},
	{"MapsLocalTaxi", icons.MapsLocalTaxi},
	{"MapsMap", icons.MapsMap},
	{"MapsMyLocation", icons.MapsMyLocation},
	{"MapsNavigation", icons.MapsNavigation},
	{"MapsNearMe", icons.MapsNearMe},
	{"MapsPersonPin", icons.MapsPersonPin},
	{"MapsPersonPinCircle", icons.MapsPersonPinCircle},
	{"MapsPinDrop", icons.MapsPinDrop},
	{"MapsPlace", icons.MapsPlace},
	{"MapsRateReview", icons.MapsRateReview},
	{"MapsRestaurant", icons.MapsRestaurant},
	{"MapsRestaurantMenu", icons.MapsRestaurantMenu},
	{"MapsSatellite", icons.MapsSatellite},
	{"MapsStoreMallDirectory", icons.MapsStoreMallDirectory},
	{"MapsStreetView", icons.MapsStreetView},
	{"MapsSubway", icons.MapsSubway},
	{"MapsTerrain", icons.MapsTerrain},
	{"MapsTraffic", icons.MapsTraffic},
	{"MapsTrain", icons.MapsTrain},
	{"MapsTram", icons.MapsTram},
	{"MapsTransferWithinAStation", icons.MapsTransferWithinAStation},
	{"MapsZoomOutMap", icons.MapsZoomOutMap},
	{"NavigationApps", icons.NavigationApps},
	{"NavigationArrowBack", icons.NavigationArrowBack},
	{"NavigationArrowDownward", icons.NavigationArrowDownward},
	{"NavigationArrowDropDown", icons.NavigationArrowDropDown},
	{"NavigationArrowDropDownCircle", icons.NavigationArrowDropDownCircle},
	{"NavigationArrowDropUp", icons.NavigationArrowDropUp},
	{"NavigationArrowForward", icons.NavigationArrowForward},
	{"NavigationArrowUpward", icons.NavigationArrowUpward},
	{"NavigationCancel", icons.NavigationCancel},
	{"NavigationCheck", icons.NavigationCheck},
	{"NavigationChevronLeft", icons.NavigationChevronLeft},
	{"NavigationChevronRight", icons.NavigationChevronRight},
	{"NavigationClose", icons.NavigationClose},
	{"NavigationExpandLess", icons.NavigationExpandLess},
	{"NavigationExpandMore", icons.NavigationExpandMore},
	{"NavigationFirstPage", icons.NavigationFirstPage},
	{"NavigationFullscreen", icons.NavigationFullscreen},
	{"NavigationFullscreenExit", icons.NavigationFullscreenExit},
	{"NavigationLastPage", icons.NavigationLastPage},
	{"NavigationMenu", icons.NavigationMenu},
	{"NavigationMoreHoriz", icons.NavigationMoreHoriz},
	{"NavigationMoreVert", icons.NavigationMoreVert},
	{"NavigationRefresh", icons.NavigationRefresh},
	{"NavigationSubdirectoryArrowLeft", icons.NavigationSubdirectoryArrowLeft},
	{"NavigationSubdirectoryArrowRight", icons.NavigationSubdirectoryArrowRight},
	{"NavigationUnfoldLess", icons.NavigationUnfoldLess},
	{"NavigationUnfoldMore", icons.NavigationUnfoldMore},
	{"NotificationADB", icons.NotificationADB},
	{"NotificationAirlineSeatFlat", icons.NotificationAirlineSeatFlat},
	{"NotificationAirlineSeatFlatAngled", icons.NotificationAirlineSeatFlatAngled},
	{"NotificationAirlineSeatIndividualSuite", icons.NotificationAirlineSeatIndividualSuite},
	{"NotificationAirlineSeatLegroomExtra", icons.NotificationAirlineSeatLegroomExtra},
	{"NotificationAirlineSeatLegroomNormal", icons.NotificationAirlineSeatLegroomNormal},
	{"NotificationAirlineSeatLegroomReduced", icons.NotificationAirlineSeatLegroomReduced},
	{"NotificationAirlineSeatReclineExtra", icons.NotificationAirlineSeatReclineExtra},
	{"NotificationAirlineSeatReclineNormal", icons.NotificationAirlineSeatReclineNormal},
	{"NotificationBluetoothAudio", icons.NotificationBluetoothAudio},
	{"NotificationConfirmationNumber", icons.NotificationConfirmationNumber},
	{"NotificationDiscFull", icons.NotificationDiscFull},
	{"NotificationDoNotDisturb", icons.NotificationDoNotDisturb},
	{"NotificationDoNotDisturbAlt", icons.NotificationDoNotDisturbAlt},
	{"NotificationDoNotDisturbOff", icons.NotificationDoNotDisturbOff},
	{"NotificationDoNotDisturbOn", icons.NotificationDoNotDisturbOn},
	{"NotificationDriveETA", icons.NotificationDriveETA},
	{"NotificationEnhancedEncryption", icons.NotificationEnhancedEncryption},
	{"NotificationEventAvailable", icons.NotificationEventAvailable},
	{"NotificationEventBusy", icons.NotificationEventBusy},
	{"NotificationEventNote", icons.NotificationEventNote},
	{"NotificationFolderSpecial", icons.NotificationFolderSpecial},
	{"NotificationLiveTV", icons.NotificationLiveTV},
	{"NotificationMMS", icons.NotificationMMS},
	{"NotificationMore", icons.NotificationMore},
	{"NotificationNetworkCheck", icons.NotificationNetworkCheck},
	{"NotificationNetworkLocked", icons.NotificationNetworkLocked},
	{"NotificationNoEncryption", icons.NotificationNoEncryption},
	{"NotificationOnDemandVideo", icons.NotificationOnDemandVideo},
	{"NotificationPersonalVideo", icons.NotificationPersonalVideo},
	{"NotificationPhoneBluetoothSpeaker", icons.NotificationPhoneBluetoothSpeaker},
	{"NotificationPhoneForwarded", icons.NotificationPhoneForwarded},
	{"NotificationPhoneInTalk", icons.NotificationPhoneInTalk},
	{"NotificationPhoneLocked", icons.NotificationPhoneLocked},
	{"NotificationPhoneMissed", icons.NotificationPhoneMissed},
	{"NotificationPhonePaused", icons.NotificationPhonePaused},
	{"NotificationPower", icons.NotificationPower},
	{"NotificationPriorityHigh", icons.NotificationPriorityHigh},
	{"NotificationRVHookup", icons.NotificationRVHookup},
	{"NotificationSDCard", icons.NotificationSDCard},
	{"NotificationSIMCardAlert", icons.NotificationSIMCardAlert},
	{"NotificationSMS", icons.NotificationSMS},
	{"NotificationSMSFailed", icons.NotificationSMSFailed},
	{"NotificationSync", icons.NotificationSync},
	{"NotificationSyncDisabled", icons.NotificationSyncDisabled},
	{"NotificationSyncProblem", icons.NotificationSyncProblem},
	{"NotificationSystemUpdate", icons.NotificationSystemUpdate},
	{"NotificationTapAndPlay", icons.NotificationTapAndPlay},
	{"NotificationTimeToLeave", icons.NotificationTimeToLeave},
	{"NotificationVibration", icons.NotificationVibration},
	{"NotificationVoiceChat", icons.NotificationVoiceChat},
	{"NotificationVPNLock", icons.NotificationVPNLock},
	{"NotificationWC", icons.NotificationWC},
	{"NotificationWiFi", icons.NotificationWiFi},
	{"PlacesACUnit", icons.PlacesACUnit},
	{"PlacesAirportShuttle", icons.PlacesAirportShuttle},
	{"PlacesAllInclusive", icons.PlacesAllInclusive},
	{"PlacesBeachAccess", icons.PlacesBeachAccess},
	{"PlacesBusinessCenter", icons.PlacesBusinessCenter},
	{"PlacesCasino", icons.PlacesCasino},
	{"PlacesChildCare", icons.PlacesChildCare},
	{"PlacesChildFriendly", icons.PlacesChildFriendly},
	{"PlacesFitnessCenter", icons.PlacesFitnessCenter},
	{"PlacesFreeBreakfast", icons.PlacesFreeBreakfast},
	{"PlacesGolfCourse", icons.PlacesGolfCourse},
	{"PlacesHotTub", icons.PlacesHotTub},
	{"PlacesKitchen", icons.PlacesKitchen},
	{"PlacesPool", icons.PlacesPool},
	{"PlacesRoomService", icons.PlacesRoomService},
	{"PlacesRVHookup", icons.PlacesRVHookup},
	{"PlacesSmokeFree", icons.PlacesSmokeFree},
	{"PlacesSmokingRooms", icons.PlacesSmokingRooms},
	{"PlacesSpa", icons.PlacesSpa},
	{"SocialCake", icons.SocialCake},
	{"SocialDomain", icons.SocialDomain},
	{"SocialGroup", icons.SocialGroup},
	{"SocialGroupAdd", icons.SocialGroupAdd},
	{"SocialLocationCity", icons.SocialLocationCity},
	{"SocialMood", icons.SocialMood},
	{"SocialMoodBad", icons.SocialMoodBad},
	{"SocialNotifications", icons.SocialNotifications},
	{"SocialNotificationsActive", icons.SocialNotificationsActive},
	{"SocialNotificationsNone", icons.SocialNotificationsNone},
	{"SocialNotificationsOff", icons.SocialNotificationsOff},
	{"SocialNotificationsPaused", icons.SocialNotificationsPaused},
	{"SocialPages", icons.SocialPages},
	{"SocialPartyMode", icons.SocialPartyMode},
	{"SocialPeople", icons.SocialPeople},
	{"SocialPeopleOutline", icons.SocialPeopleOutline},
	{"SocialPerson", icons.SocialPerson},
	{"SocialPersonAdd", icons.SocialPersonAdd},
	{"SocialPersonOutline", icons.SocialPersonOutline},
	{"SocialPlusOne", icons.SocialPlusOne},
	{"SocialPoll", icons.SocialPoll},
	{"SocialPublic", icons.SocialPublic},
	{"SocialSchool", icons.SocialSchool},
	{"SocialSentimentDissatisfied", icons.SocialSentimentDissatisfied},
	{"SocialSentimentNeutral", icons.SocialSentimentNeutral},
	{"SocialSentimentSatisfied", icons.SocialSentimentSatisfied},
	{"SocialSentimentVeryDissatisfied", icons.SocialSentimentVeryDissatisfied},
	{"SocialSentimentVerySatisfied", icons.SocialSentimentVerySatisfied},
	{"SocialShare", icons.SocialShare},
	{"SocialWhatsHot", icons.SocialWhatsHot},
	{"ToggleCheckBox", icons.ToggleCheckBox},
	{"ToggleCheckBoxOutlineBlank", icons.ToggleCheckBoxOutlineBlank},
	{"ToggleIndeterminateCheckBox", icons.ToggleIndeterminateCheckBox},
	{"ToggleRadioButtonChecked", icons.ToggleRadioButtonChecked},
	{"ToggleRadioButtonUnchecked", icons.ToggleRadioButtonUnchecked},
	{"ToggleStar", icons.ToggleStar},
	{"ToggleStarBorder", icons.ToggleStarBorder},
	{"ToggleStarHalf", icons.ToggleStarHalf},
}
